package main

import (
	"github.com/gogo/protobuf/vanity/command"
	"github.com/jackskj/protoc-gen-map/plugin"
)

func main() {
	plugin, _ := plugin.New()
	response := command.GeneratePlugin(
		command.Read(),
		plugin,
		".pb.map.go",
	)
	command.Write(response)
}
