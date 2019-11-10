package templates

var server = `
type {{ .ServiceName}}MapServer struct {
	DB  *sql.DB
	mapperGenMux sync.Mutex

	{{ range $key, $val := .MapperNames }}
	{{ $key }}Mapper *mapper.Mapper
	{{ $key }}Callbacks {{ $.ServiceName}}{{ $key }}Callbacks
	{{- end }}
}
`
