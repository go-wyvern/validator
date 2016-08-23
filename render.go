package validator

import (
	"github.com/go-wyvern/Leego"
	"io"
	"strings"
	"html/template"
	"os"
	"fmt"
	"encoding/json"
)

var AppApis []Api

const codeTag = "```"

type Api struct {
	Description string
	Method      string
	Path        string
	Handler   leego.HandlerFunc
	SuccessStdOut interface{}
	SuccessFormat []byte
	FailStdOut interface{}
	FailFormat  []byte
	StdFormat   string
	CodeTag     string
	Validator *Validator
}

type Module struct {
	ModuleName string
	Apis   []Api
}

type Project struct {
	ProjectName string
	Modules      []Module
}

func Find(method, path string) *Validator {
	for _,r:=range AppApis {
		if r.Method==method&& r.Path==path {
			return r.Validator
		}
	}
	return nil
}

func NewProject(name string) *Project {
	p := new(Project)
	p.ProjectName = name
	return p
}

func NewModule(module_name string) *Module {
	module := new(Module)
	module.ModuleName = module_name
	return module
}

func NewApi(method, path,d string, h leego.HandlerFunc, v *Validator) *Api {
	r := new(Api)
	r.Method = method
	r.Description=d
	r.Path = path
	r.Handler = h
	r.Validator = v
	r.CodeTag=codeTag
	return r
}

func (c *Api) SetSuccessStdOut(s interface{}) {
	c.SuccessStdOut = s
	if c.StdFormat=="json" {
		c.SuccessFormat,_=json.MarshalIndent(s,"","  ")
	}
}

func (c *Api) SetFailStdOut(s interface{}) {
	c.FailStdOut = s
	if c.StdFormat=="json" {
		c.FailFormat,_=json.MarshalIndent(s,"","  ")
	}
}

func (c *Module) Use(a Api) *Module {
	c.Apis = append(c.Apis, a)
	AppApis = append(AppApis, a)
	return c
}

func (c *Project) Use(m Module) *Project {
	c.Modules = append(c.Modules, m)
	return c
}

func (c *Project) RenderMarkdown(filename string,app *Project) error {
	var err error
	f,err:=os.Create(filename)
	if err!=nil{
		fmt.Println(err.Error())
		return err
	}
	err=tmpl(f,MarkdownTemplate,app)
	if err != nil {
		return err
	}
	return nil
}

var MarkdownTemplate =`{{with .}}# {{.ProjectName}}
{{range .Modules}}
## {{.ModuleName}}
{{range .Apis}}
### {{.Method}} {{.Path}} {{.Description}}

请求参数:

| 名称 | 类型 | 说明           | 是否必须  |
| -----|:-----:|:---------:|:-----:|
{{range $name, $params :=.Validator.ApiParams}}|**{{$name}}**|{{$params.Type}}|{{$params.Description}}|{{$params.Require}}|
{{end}}
请求正确返回:

{{.CodeTag}}
{{.SuccessFormat|printf "%s"|unescaped}}
{{.CodeTag}}

请求错误返回:

{{.CodeTag}}
{{.FailFormat | printf "%s"|unescaped}}
{{.CodeTag}}

{{end}}{{end}}{{end}}
`

func tmpl(w io.Writer, text string, data interface{})error {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": func(s template.HTML) template.HTML {
		return template.HTML(strings.TrimSpace(string(s)))
	}})

	t.Funcs(template.FuncMap{"unescaped": func(x string) interface{} {
		return template.HTML(x)
	}})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		return err
	}
	return nil
}