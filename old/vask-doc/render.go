package main

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/drgo/mdson"
)

func runTemplate(env *Env, templateFileName, dataFileName string) error {
	templateName := strings.TrimSuffix(templateFileName, filepath.Ext(templateFileName))
	destFileName := "/tmp/" + strings.TrimSuffix(dataFileName, filepath.Ext(dataFileName)) + ".html"
	templateFileName = env.GetSourceFilePath(templateFileName)
	dataFileName = env.GetSourceFilePath(dataFileName)

	t, err := template.New("tt").Funcs(funcList).ParseFiles(templateFileName)
	if err != nil {
		return fmt.Errorf("failed to parse file '%s': %v", templateFileName, err)
	}
	data, err := parse(dataFileName)
	// fmt.Println(data)
	if err != nil {
		return fmt.Errorf("failed to parse file '%s': %v", dataFileName, err)
	}

	out, err := env.CreateFile(destFileName)
	if err != nil {
		return fmt.Errorf("failed to create output file '%s': %v", destFileName, err)
	}
	defer out.Close()
	err = t.ExecuteTemplate(out, templateName, data)
	if err != nil {
		return fmt.Errorf("failed to render file '%s': %v", dataFileName, err)
	}
	return nil
}

func parse(fileName string) (mdson.Node, error) {
	root, err := mdson.ParseFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file '%s': %v", fileName, err)
	}
	if root.Kind() != "Block" {
		return nil, fmt.Errorf("parser returned unexpected type: root is not a block")
	}
	return root, nil
}

// Template funcs
func toJS(s string) template.JS {
	return template.JS(s)
}

func toHTML(s string) template.HTML {
	return template.HTML(s)
}

var funcList = template.FuncMap{
	"title":  strings.Title,
	"toJS":   toJS,
	"toHTML": toHTML,
}
