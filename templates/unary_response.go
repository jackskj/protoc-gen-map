package templates

var unaryResponse = `
func (m *{{ .ServiceName }}MapServer) {{ .MethodName }}(ctx context.Context, r *{{ .RequestName }}) (*{{ .ResponseName }}, error) {
	sqlBuffer := &bytes.Buffer{}
	if err := {{ .SqlTemplateName }}.ExecuteTemplate(sqlBuffer, "{{.MethodName }}", r); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rawSql := sqlBuffer.String()
	{{ if eq .QueryType "Exec" }}
	_, err := m.DB.Exec(rawSql)
	if err != nil {
		log.Printf("error executing query.\n {{ .RequestName }} request: %s \n,query: %s \n error: %s", r, rawSql, err)
		return nil, status.Error(codes.InvalidArgument, "request generated malformed query")
	}
	resp :={{ .ResponseName }}{}
        return &resp, nil
	{{ else if eq .QueryType "Query" }}
	rows, err := m.DB.Query(rawSql)
	defer rows.Close()
	if err != nil {
		log.Printf("error executing query.\n {{ .RequestName }} request: %s \n,query: %s \n error: %s", r, rawSql, err)
		return nil, status.Error(codes.InvalidArgument, "request generated malformed query")
	}
	if m.{{ .MapperName }}Mapper == nil {
		m.mapperGenMux.Lock()
		m.{{ .MapperName }}Mapper, err =  mapper.New(rows, &{{ .ResponseName }}{})
		m.mapperGenMux.Unlock()
		if err != nil {
			log.Printf("error generating {{ .MapperName }}Mapper: %s", err)
			return nil, status.Error(codes.Internal, "error generating {{ .ResponseName }} mapping")
		}
		m.{{ .MapperName }}Mapper.Log()
	}
	respMap := m.{{ .MapperName }}Mapper.NewResponseMapping()
	if err := m.{{ .MapperName }}Mapper.GetValues(rows, respMap); err != nil{
		log.Printf("error loading data for {{ .MethodName }}: %s", err)
		return nil, status.Error(codes.Internal, "error loading data")
	}
	if err := m.{{ .MapperName }}Mapper.MapResponse(respMap); err != nil{
		log.Printf("error mappig {{ .MapperName }}Mapper: %s", err)
		m.{{ .MapperName }}Mapper.Error = nil
		return nil, status.Error(codes.Internal, "error mappig {{ .ResponseName }}")
	}
	m.{{ .MapperName }}Mapper.Log()
	if len(respMap.Responses) == 0{
		//No Responses found
		return new({{ .ResponseName }}), nil
	} else {
		return respMap.Responses[0].(*{{ .ResponseName }}) , nil
	}
	{{ end }}
}
`
