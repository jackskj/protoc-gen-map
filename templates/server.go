package templates

var server = `
type {{ .ServiceName}}MapServer struct {
	DB  *sql.DB
	{{- range $key, $val := .MapperNames }}
	{{ $key }}Mapper *mapper.Mapper
	{{- end }}

	mapperGenMux sync.Mutex
}
`
