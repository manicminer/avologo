package main

import (
	"html/template"
	"io"
	"github.com/labstack/echo"
)

type TemplateRenderer struct {
	templates *template.Template
}

/*
	Renders a template document
*/
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}