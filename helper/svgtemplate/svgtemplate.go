// Package svgtemplate embeds the template files into the build go binary.
package svgtemplate

import (
	"embed"
	"html/template"
)

//go:embed templates/*.gosvg
var embedFS embed.FS

// GetSVGTemplates - parse and get all SVG templates
func GetSVGTemplates() (*template.Template, error) {
	svgTemplates, err := template.ParseFS(embedFS, "templates/*.gosvg")
	if err != nil {
		return nil, err
	}
	return svgTemplates, nil
}
