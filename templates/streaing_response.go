package templates

var streamingResponse = `
func (m *{{ .ServiceName }}MapServer) {{ .MethodName }}(r *{{ .RequestName }}, stream {{ .ServiceName }}_{{ .MethodName }}Server) error {
	sqlBuffer := &bytes.Buffer{}
	if err := {{ .SqlTemplateName }}.ExecuteTemplate(sqlBuffer, "{{.MethodName }}", r); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	rawSql := sqlBuffer.String()
	rows, err := m.DB.Query(rawSql)
	defer rows.Close()
	if err != nil {
		log.Printf("error executing query.\n {{ .RequestName }} request: %s \n,query: %s \n error: %s", r, rawSql, err)
		return status.Error(codes.InvalidArgument, "request generated malformed query")
	}
	if m.{{ .MapperName }}Mapper == nil {
		m.mapperGenMux.Lock()
		m.{{ .MapperName }}Mapper, err =  mapper.New(rows, &{{ .ResponseName }}{})
		m.mapperGenMux.Unlock()
		if err != nil {
			log.Printf("error generating {{ .MapperName }}Mapper: %s", err)
			return status.Error(codes.Internal, "error generating {{ .ResponseName }} mapping")
		}
		m.{{ .MapperName }}Mapper.Log()
	}
	respMap := m.{{ .MapperName }}Mapper.NewResponseMapping()
	if err := m.{{ .MapperName }}Mapper.GetValues(rows, respMap); err != nil{
		log.Printf("error loading data for {{ .MethodName }}: %s", err)
		return status.Error(codes.Internal, "error loading data")
	}
	if err := m.{{ .MapperName }}Mapper.MapResponse(respMap); err != nil{
		log.Printf("error mappig {{ .MapperName }}Mapper: %s", err)
		m.{{ .MapperName }}Mapper.Error = nil
		return status.Error(codes.Internal, "error mappig {{ .ResponseName }}")
	}
	m.{{ .MapperName }}Mapper.Log()
	for _, resp := range respMap.Responses {
		if err := stream.Send(resp.(*{{ .ResponseName }})); err != nil {
			return err
		}
	}
	return nil
}
`
