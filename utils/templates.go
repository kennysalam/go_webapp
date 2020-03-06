package utils

import (
	"html/template"
	"net/http"
)

//Templates html templates
var templates *template.Template

//LoadTemplates load html templates
func LoadTemplates(pattern string) {
	templates = template.Must(template.ParseGlob(pattern))
}

//ExecuteTemplate execute template
func ExecuteTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl, data)
}
