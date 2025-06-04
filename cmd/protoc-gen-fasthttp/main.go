// main.go
package main

import (
	_ "embed"
	"log"
	"net/http"
	"text/template"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

type RpcData struct {
	RPCName      string
	HTTPVerb     string
	HTTPPath     string
	RequestType  string
	ResponseType string
	HandlerName  string
}

type TemplateData struct {
	Package       string
	ServiceName   string
	RPCs          []RpcData
	ImportPath    string
	GoPackagePath string
	ServerName    string
}

type HttpData struct {
	Verb string
	Path string
}

func extractHTTPInfo(m *protogen.Method) (hd HttpData, ok bool) {
	opts := m.Desc.Options()
	if opts == nil {
		return
	}

	ext := proto.GetExtension(opts, annotations.E_Http)
	rule, ok := ext.(*annotations.HttpRule)
	if !ok || rule == nil {
		return
	}

	switch pattern := rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		hd.Verb = http.MethodGet
		hd.Path = pattern.Get
	case *annotations.HttpRule_Post:
		hd.Verb = http.MethodPost
		hd.Path = pattern.Post
	case *annotations.HttpRule_Put:
		hd.Verb = http.MethodPut
		hd.Path = pattern.Put
	case *annotations.HttpRule_Delete:
		hd.Verb = http.MethodDelete
		hd.Path = pattern.Delete
	case *annotations.HttpRule_Patch:
		hd.Verb = http.MethodPatch
		hd.Path = pattern.Patch
	}

	ok = true

	return
}

//go:embed handler.tmpl
var tplContent string

func main() {
	tmpl := template.Must(template.New("handler").Parse(tplContent))

	protogen.Options{}.Run(func(plugin *protogen.Plugin) error {
		for _, file := range plugin.Files {
			if !file.Generate {
				continue
			}
			for _, service := range file.Services {
				var rpcs []RpcData
				for _, method := range service.Methods {
					hd, ok := extractHTTPInfo(method)
					if !ok {
						continue
					}

					rpcs = append(rpcs, RpcData{
						RPCName:      method.GoName,
						HTTPVerb:     hd.Verb,
						HTTPPath:     hd.Path,
						RequestType:  method.Input.GoIdent.GoName,
						ResponseType: method.Output.GoIdent.GoName,
						HandlerName:  service.GoName + method.GoName,
					})
				}

				if len(rpcs) == 0 {
					continue
				}

				filename := file.GeneratedFilenamePrefix + ".pb.fh.go"
				g := plugin.NewGeneratedFile(filename, file.GoImportPath)

				data := TemplateData{
					Package:       string(file.GoPackageName),
					ServiceName:   service.GoName,
					ServerName:    service.GoName + "ServerFastHttp",
					ImportPath:    string(file.GoImportPath),
					RPCs:          rpcs,
					GoPackagePath: string(file.GoImportPath),
				}

				if err := tmpl.Execute(g, data); err != nil {
					log.Fatalf("failed to render template: %v", err)
				}
			}
		}
		return nil
	})
}
