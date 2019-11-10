package templates

var unaryResponse = `

type {{ .ServiceName}}{{ .MethodName }}Callbacks struct {
	BeforeQueryCallback []func(queryString string, req *{{ .RequestName }}) error
	AfterQueryCallback []func(queryString string, req *{{ .RequestName }}, resp *{{ .ResponseName }} ) error
	Cache func(queryString string, req *{{ .RequestName }}) (*{{ .ResponseName }}, error)
}

func (m *{{ .ServiceName }}MapServer) Register{{ .MethodName }}BeforeQueryCallback(callbacks ...func(queryString string, req *{{ .RequestName }}) error ) {
	for _, callback := range callbacks {
		m.{{ .MethodName }}Callbacks.BeforeQueryCallback = append(m.{{ .MethodName }}Callbacks.BeforeQueryCallback, callback)
	}
}

func (m *{{ .ServiceName }}MapServer) Register{{ .MethodName }}AfterQueryCallback(callbacks ...func(queryString string, req *{{ .RequestName }}, resp *{{ .ResponseName }} ) error ) {
	for _, callback := range callbacks {
		m.{{ .MethodName }}Callbacks.AfterQueryCallback = append(m.{{ .MethodName }}Callbacks.AfterQueryCallback, callback)
	}
}

func (m *{{ .ServiceName }}MapServer) Register{{ .MethodName }}Cache(cache func(queryString string, req *{{ .RequestName }}) (*{{ .ResponseName }}, error )) {
	m.{{ .MethodName }}Callbacks.Cache = cache
}

func (m *{{ .ServiceName }}MapServer) {{ .MethodName }}(ctx context.Context, r *{{ .RequestName }}) (*{{ .ResponseName }}, error) {
	sqlBuffer := &bytes.Buffer{}
	if err := {{ .SqlTemplateName }}.ExecuteTemplate(sqlBuffer, "{{.MethodName }}", r); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rawSql := sqlBuffer.String()
	for _, callback := range m.{{ .MethodName }}Callbacks.BeforeQueryCallback {
		if err := callback(rawSql, r); err != nil {
			log.Println(err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	if  m.{{ .MethodName }}Callbacks.Cache != nil {
		if resp, err :=  m.{{ .MethodName }}Callbacks.Cache(rawSql, r); err == nil {
			if resp != nil {
				return resp, nil
			}
		} else {
			log.Println(err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	{{ if eq .QueryType "Exec" }}
	_, err := m.DB.Exec(rawSql)
	if err != nil {
		log.Printf("error executing query.\n {{ .RequestName }} request: %s \n,query: %s \n error: %s", r, rawSql, err)
		return nil, status.Error(codes.InvalidArgument, "request generated malformed query")
	}
	for _, callback := range m.{{ .MethodName }}Callbacks.AfterQueryCallback {
		if err := callback(rawSql, r, nil); err != nil {
			log.Println(err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	resp :={{ .ResponseName }}{}
        return &resp, nil
	{{ else if eq .QueryType "Query" }}
	rows, err := m.DB.Query(rawSql)
	if err != nil {
		log.Printf("error executing query.\n {{ .RequestName }} request: %s \n,query: %s \n error: %s", r, rawSql, err)
		return nil, status.Error(codes.InvalidArgument, "request generated malformed query")
	} else {
		defer rows.Close()
	}
	if m.{{ .MapperName }}Mapper == nil {
		m.mapperGenMux.Lock()
		m.{{ .MapperName }}Mapper, err =  mapper.New("{{ .MethodName }}", rows, &{{ .ResponseName }}{})
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
	var response  *{{ .ResponseName }}
	if len(respMap.Responses) == 0{
		//No Responses found
		response = &{{ .ResponseName }}{}
	} else {
		response = respMap.Responses[0].(*{{ .ResponseName }})
	}
	for _, callback := range m.{{ .MethodName }}Callbacks.AfterQueryCallback {
		if err := callback(rawSql, r, response); err != nil {
			log.Println(err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	m.{{ .MapperName }}Mapper.Log()
	return response, nil
	{{ end }}
}
`
