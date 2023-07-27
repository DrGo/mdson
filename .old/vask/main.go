package main

import (
	"fmt"
	"html/template"
	"os"
)


func check(err error) {
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func main(){
  title := "home"

    data := map[string]interface{}{
        "title":  title,
        "header": "My Header",
        "footer": "My Footer",
    }

    check(templates.ExecuteTemplate(os.Stdout, "homeHTML", data))
    
}

var templates *template.Template

func getTemplates(dir string) (templates *template.Template, err error) {
    return  template.ParseGlob(dir+"/*")  
}

func init() {
    templates, _ = getTemplates("site/layout")
}
