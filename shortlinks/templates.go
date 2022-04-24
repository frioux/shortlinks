package shortlinks

import (
	"embed"
	"html/template"
)

//go:embed templates/*
var templates embed.FS

var tpl *template.Template

func init() {
	var err error
	tpl, err = template.ParseFS(templates, "templates/*")
	if err != nil {
		panic(err)
	}
}
