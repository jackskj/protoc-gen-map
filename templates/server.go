package templates

var server = `
type {{ .ServiceName}}MapServer struct {
	DB  *sql.DB
	Dialect string
	{{ range $key, $val := .MapperNames }}
	{{ $key }}Mapper *mapper.Mapper
	{{ $key }}Callbacks {{ $.ServiceName}}{{ $key }}Callbacks
	{{- end }}

	mapperGenMux sync.Mutex
}
`
