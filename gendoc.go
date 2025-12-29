package swagger

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"text/template"
)

//go:embed doc_template.md
var markdownTlp string

//go:embed doc_template.html
var htmlTlp string

func GenerateMarkdown(api []*Api) string {
	return generateFromTemplate(api, markdownTlp)
}

func GenerateHtml(api []*Api) string {
	return generateFromTemplate(api, htmlTlp)
}

func generateFromTemplate(api []*Api, tlp string) string {

	apidocs := []*apiDoc{}
	for _, a := range api {
		apidocs = append(apidocs, &apiDoc{
			Api: a,
			Req: a.RequestSchema.Doc(),
			Res: a.ResponseSchema.Doc(),
			ReqExample: func() string {
				body := a.RequestSchema.GenExampleJson()
				if body == "{}" {
					body = ""
				}
				query := strings.Join(a.RequestSchema.genExampleQuery(), "&")
				if query != "" {
					query = "?" + query
				}
				return fmt.Sprintf("%s %s%s\n\n", a.Method, a.RequestSchema.generateExamplePath(a.Route), query) + body
			}(),
			ResExample: a.ResponseSchema.GenExampleJson(),
		})
	}
	t, err := template.New("swagger").Parse(tlp)
	if err != nil {
		panic(err)
	}
	bf := &bytes.Buffer{}
	err = t.Execute(bf, apidocs)
	if err != nil {
		panic(err)
	}
	return bf.String()
}

type apiDoc struct {
	Id         string
	Api        *Api
	ReqExample any
	ResExample any
	Req        []*FiledDoc
	Res        []*FiledDoc
}

type FiledDoc struct {
	Field       string
	Type        string
	Enum        string
	Required    bool
	Description string
	Location    string
	Default     string
	Binding     string
}
