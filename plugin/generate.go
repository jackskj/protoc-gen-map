package plugin

import (
	"bytes"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	protogen "github.com/golang/protobuf/protoc-gen-go/generator"
	"strings"
)

//RPC Method, corresponding to a single sql statement
type Server struct {
	ServiceName string
	MapperNames map[string]bool
}

//RPC Method, corresponding to a single sql statement
type Method struct {
	MethodName      string
	ServiceName     string
	MapperName      string
	RequestName     string
	ResponseName    string
	SqlTemplateName string
	QueryType       string
}

// used to store all messages in case of proto imports
type Message struct {
	Package       string
	GoPackageName string
	GoType        string
}

// Generates mapper file
func (p *SqlPlugin) Generate(file *generator.FileDescriptor) {
	p.printHeader(file)
	p.generateEnumVals(file)
	p.generateMethods(file)
	p.PrintSQLTemplates(file)
}

// Generates method for each rpc
func (p *SqlPlugin) generateMethods(file *generator.FileDescriptor) {
	methodBuff := bytes.Buffer{}

	for _, service := range file.GetService() {
		p.generateServiceServer(service, file)
		for _, method := range service.GetMethod() {
			if _, found := p.matchingTemplates[method.GetName()]; !found {
				continue
			}
			if method.GetClientStreaming() {
				p.Error(method.GetName() + " is client streaming, protoc-gen-map does not allow client streaming")
			}
			p.method = &Method{
				MethodName:      method.GetName(),
				ServiceName:     service.GetName(),
				MapperName:      method.GetName(),
				RequestName:     p.getMessageName(method.GetInputType(), file),
				ResponseName:    p.getMessageName(method.GetOutputType(), file),
				SqlTemplateName: p.SqlTemplateName,
				QueryType:       p.getQueryType(method),
			}
			if method.GetServerStreaming() == true {
				err := p.genTemplate.ExecuteTemplate(&methodBuff, "streamingResponse", p.method)
				if err != nil {
					p.Error(err.Error())
				}
				p.setStreamingImports()
			} else {
				err := p.genTemplate.ExecuteTemplate(&methodBuff, "unaryResponse", p.method)
				if err != nil {
					p.Error(err.Error())
				}
				p.setUnaryImports()
			}
			p.P(methodBuff.String())
			methodBuff.Reset()
		}
	}
}

// Generates the service server struct named <service name>MapServer, which implements the
// corresponging interface generated by go proto complier
func (p *SqlPlugin) generateServiceServer(service *descriptor.ServiceDescriptorProto, file *generator.FileDescriptor) {
	serverBuff := bytes.Buffer{}
	mapperNames := make(map[string]bool)
	p.server = &Server{
		ServiceName: service.GetName(),
		MapperNames: mapperNames,
	}
	for _, method := range service.GetMethod() {
		p.server.MapperNames[method.GetName()] = true
	}
	if len(p.server.MapperNames) != 0 {
		p.Pkg["mapper"] = true
	}
	err := p.genTemplate.ExecuteTemplate(&serverBuff, "server", p.server)

	if err != nil {
		p.Error(err.Error())
	}
	p.P(serverBuff.String())
	serverBuff.Reset()
	p.Pkg["sync"] = true
	p.Pkg["sql"] = true
}

func (p *SqlPlugin) getQueryType(method *descriptor.MethodDescriptorProto) (queryType string) {
	queryType = "Query"
	name := method.GetName()
	outSlice := strings.Split(method.GetOutputType(), ".")
	out := outSlice[len(outSlice)-1]
	for prefix, _ := range execOperations {
		if strings.HasPrefix(strings.ToLower(name), prefix) {
			queryType = "Exec"
		}
	}
	for prefix, _ := range execResponses {
		if strings.HasPrefix(strings.ToLower(out), prefix) {
			queryType = "Exec"
		}
	}
	if queryType == "Exec" && method.GetClientStreaming() {
		p.Error(method.GetName() +
			" does not expect response rows and is server streaming streaming.\n" +
			"Please read important notes at https://github.com/jackskj/protoc-gen-map#important-notes",
		)
	}
	return
}

// Requests and responses to rpc can be imported, however, method descriptor only returns
// message type names, this function finds thether the messages are imported and concatinates it
// with correct go package prefix if nececary
// Also I take into account that protoc makes message names camel case
func (p *SqlPlugin) getMessageName(msg string, file *generator.FileDescriptor) string {
	name_slice := strings.Split(msg, ".")

	//if not imported
	if strings.Join(name_slice[0:len(name_slice)-1], ".") == "."+file.GetPackage() {
		return strings.Title(name_slice[len(name_slice)-1])
	} else {
		packageName := strings.Join(name_slice[1:len(name_slice)-1], ".")
		typeName := name_slice[len(name_slice)-1]
		var goPackageName string
		for _, file := range p.Generator.AllFiles().GetFile() {
			if file.GetPackage() == packageName {
				if file.GetOptions().GetGoPackage() != "" {
					goPackageSplit := strings.Split(file.GetOptions().GetGoPackage(), "/")
					//in case go import does not match directory name
					goAliasSplit := strings.Split(goPackageSplit[len(goPackageSplit)-1], ";")
					goPackageName = goAliasSplit[len(goAliasSplit)-1]
					imports[alias(goPackageName)] = strings.Join(goPackageSplit[0:len(goPackageSplit)-1], "/") +
						"/" + goAliasSplit[0]
					p.Pkg[alias(goPackageName)] = true
				} else {
					goPackageSplit := strings.Split(file.GetPackage(), ".")
					goPackageName = goPackageSplit[len(goPackageSplit)-1]
					imports[alias(goPackageName)] = strings.Join(goPackageSplit, "/")
					p.Pkg[alias(goPackageName)] = true
				}
				break
			}
		}
		return goPackageName + "." + strings.Title(typeName)
	}
}

func (p *SqlPlugin) generateEnumVals(file *generator.FileDescriptor) {
	p.EnumValueMaps = make(map[string]map[string]int32)
	p.findEnums(file.GetEnumType(), []string{})
	for _, message := range file.GetMessageType() {
		p.findEnums(message.GetEnumType(), []string{message.GetName()})
		p.findNestedEnums(message, []string{message.GetName()})
	}
	if len(p.EnumValueMaps) > 0 {
		p.setEnumImports()
		buff := bytes.Buffer{}
		err := p.genTemplate.ExecuteTemplate(&buff, "enumValueMaps", p)
		if err != nil {
			p.Error(err.Error())
		}
		err = p.genTemplate.ExecuteTemplate(&buff, "initFunc", nil)
		if err != nil {
			p.Error(err.Error())
		}
		p.P(buff.String())
	}
}

func (p *SqlPlugin) findNestedEnums(message *descriptor.DescriptorProto, typeName []string) {
	for _, nestedMsg := range message.GetNestedType() {
		nestedTypeName := append(typeName, nestedMsg.GetName())
		p.findEnums(nestedMsg.GetEnumType(), nestedTypeName)
		p.findNestedEnums(nestedMsg, nestedTypeName)
	}
}

func (p *SqlPlugin) findEnums(enums []*descriptor.EnumDescriptorProto, typeName []string) {
	for _, enum := range enums {
		enumValueMap := make(map[string]int32)
		for _, value := range enum.GetValue() {
			enumValueMap[value.GetName()] = value.GetNumber()
		}
		ccTypeName := protogen.CamelCase(strings.Join(append(typeName, enum.GetName()), "_"))
		p.EnumValueMaps[ccTypeName] = enumValueMap
	}
}

func (p *SqlPlugin) printHeader(file *generator.FileDescriptor) {
	p.P("// Code generated by protoc-gen-map. DO NOT EDIT.")
	p.P("// To Use:")
	p.P("// 1. Instantiate MapperServers with sql.DB connection")
	p.P("// 2. Register MapperServer as the gRPC service server")
	p.P("// 3. Begin serving")
}

var execOperations = map[string]bool{
	"delete": true,
	"insert": true,
	"update": true,
	"create": true,
}
var execResponses = map[string]bool{
	"empty": true,
	"nil":   true,
	"null":  true,
}
