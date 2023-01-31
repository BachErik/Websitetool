package main

import (
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"text/template"
)

var allTemplates []string

func initializeTemplates() {
	files, err := ioutil.ReadDir("./views/partials")
	if err != nil {
		panic(err.Error())
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".tmpl") {
			allTemplates = append(allTemplates, "./views/partials/"+filename)
		}
	}
}

func serveTemplate(responseWriter http.ResponseWriter, file string, functions template.FuncMap, data interface{}) {
	var files []string
	files = append(files, file)
	files = append(files, allTemplates...)
	name := path.Base(files[0])
	templates := template.New(name)
	if functions != nil {
		templates.Funcs(functions)
	}
	templates, err := templates.ParseFiles(files...)
	if err != nil {
		panic(err.Error())
	}
	err = templates.Execute(responseWriter, data)
	if err != nil {
		panic(err.Error())
	}
}
