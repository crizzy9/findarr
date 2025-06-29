package web

import (
	"html/template"
	"path/filepath"
)

// TemplateManager handles loading and caching of templates
// Usage: templates := web.NewTemplateManager([list of template files])
//
// Example:
//  mgr := web.NewTemplateManager([]string{"web/templates/base.html", "web/templates/index.html"})
//  err := mgr.Load()
//  tpl := mgr.Get("index")

type TemplateManager struct {
	files     []string
	templates map[string]*template.Template
}

func NewTemplateManager(files []string) *TemplateManager {
	return &TemplateManager{
		files:     files,
		templates: make(map[string]*template.Template),
	}
}

// Load parses and caches templates
func (tm *TemplateManager) Load() error {
	// First, find the base template (assuming it's named base.html)
	var baseFile string
	for _, file := range tm.files {
		if filepath.Base(file) == "base.html" {
			baseFile = file
			break
		}
	}

	// Parse each template with the base template
	for _, file := range tm.files {
		name := filepath.Base(file)
		if name == "base.html" {
			// Skip base template as standalone
			continue
		}

		// Parse with base template to support template inheritance
		tpl, err := template.ParseFiles(baseFile, file)
		if err != nil {
			return err
		}
		tm.templates[name] = tpl
	}

	// Also add the base template itself
	baseTpl, err := template.ParseFiles(baseFile)
	if err != nil {
		return err
	}
	tm.templates["base.html"] = baseTpl

	return nil
}

// Get returns a parsed template by filename
func (tm *TemplateManager) Get(name string) *template.Template {
	return tm.templates[name]
}

