package templates

var streamingResponse = `

type {{ .ServiceName}}{{ .MethodName }}Callbacks struct {
	BeforeQueryCallback []func(queryString string, req *{{ .RequestName }}) error
	AfterQueryCallback []func(queryString string, req *{{ .RequestName }}, resp []*{{ .ResponseName }} ) error
	Cache func(queryString string, req *{{ .RequestName }}) ([]*{{ .ResponseName }}, error)
}

func (m *{{ .ServiceName }}MapServer) Register{{ .MethodName }}BeforeQueryCallback(callbacks ...func(queryString string, req *{{ .RequestName }}) error ) {
	for _, callback := range callbacks {
		m.{{ .MethodName }}Callbacks.BeforeQueryCallback = append(m.{{ .MethodName }}Callbacks.BeforeQueryCallback, callback)

	}
}

func (m *{{ .ServiceName }}MapServer) Register{{ .MethodName }}AfterQueryCallback(callbacks ...func(queryString string, req *{{ .RequestName }}, resp []*{{ .ResponseName }} ) error ) {
	for _, callback := range callbacks {
		m.{{ .MethodName }}Callbacks.AfterQueryCallback = append(m.{{ .MethodName }}Callbacks.AfterQueryCallback, callback)
	}
}

func (m *{{ .ServiceName }}MapServer) Register{{ .MethodName }}Cache(cache func(queryString string, req *{{ .RequestName }}) ([]*{{ .ResponseName }}, error )) {
	m.{{ .MethodName }}Callbacks.Cache = cache
}

func (m *{{ .ServiceName }}MapServer) {{ .MethodName }}(r *{{ .RequestName }}, stream {{ .ServiceName }}_{{ .MethodName }}Server) error {
	sqlBuffer := &bytes.Buffer{}
	if err := {{ .SqlTemplateName }}.ExecuteTemplate(sqlBuffer, "{{.MethodName }}", r); err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	rawSql := sqlBuffer.String()
	for _, callback := range m.{{ .MethodName }}Callbacks.BeforeQueryCallback {
		if err := callback(rawSql, r); err != nil {
			log.Println(err.Error())
			return status.Error(codes.Internal, err.Error())
		}
	}
	if  m.{{ .MethodName }}Callbacks.Cache != nil {
		if responses, err :=  m.{{ .MethodName }}Callbacks.Cache(rawSql, r); err == nil {
			if responses != nil {
				for _, resp := range responses {
					if err := stream.Send(resp); err != nil {
						return status.Error(codes.Internal, err.Error())
					}
				}
				return nil
			}
		} else {
			log.Println(err.Error())
			return status.Error(codes.Internal, err.Error())
		}
	}
	preparedSql, args, err := mapper.PrepareQuery(m.Dialect, sqlBuffer.Bytes())
	if err != nil {
		log.Printf("error preparing sql query.\n {{ .RequestName }} request: %s \n query: %s \n error: %s", r, rawSql, err)
		return status.Error(codes.InvalidArgument, "Request generated malformed query.")
	}
	rows, err := m.DB.QueryContext(stream.Context(), preparedSql, args...)
	if stream.Context().Err() == context.Canceled {
		return status.Error(codes.Canceled, "Client cancelled.")
	} else if err != nil {
		log.Printf("error executing query.\n {{ .RequestName }} request: %s \n query: %s \n error: %s", r, rawSql, err)
		return status.Error(codes.Internal, err.Error())
	} else {
		defer rows.Close()
	}
	if m.{{ .MapperName }}Mapper == nil {
		m.mapperGenMux.Lock()
		m.{{ .MapperName }}Mapper, err =  mapper.New("{{ .MethodName }}", rows, &{{ .ResponseName }}{})
		m.mapperGenMux.Unlock()
		if err != nil {
			log.Printf("Error generating {{ .MapperName }}Mapper: %s", err)
			return status.Error(codes.Internal, "Error generating {{ .ResponseName }} mapping.")
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
	var responses []*{{ .ResponseName }}
	for _, resp := range respMap.Responses {
		responses = append(responses, resp.(*{{ .ResponseName }}))
	}
	for _, callback := range m.{{ .MethodName }}Callbacks.AfterQueryCallback {
		if err := callback(rawSql, r, responses); err != nil {
			log.Println(err.Error())
			return status.Error(codes.Internal, err.Error())
		}
	}
	m.{{ .MapperName }}Mapper.Log()
	for _, resp := range responses {
		if err := stream.Send(resp); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
	return nil
}
`
