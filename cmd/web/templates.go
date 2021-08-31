package main

import (
	"html/template"
	"path/filepath"

	"zapmal/snippetbox/pkg/models"
)

type TemplateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache(directory string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(directory, "*.page.tmpl"))

	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		templateSet, err := template.ParseFiles(page)

		if err != nil {
			return nil, err
		}

		templateSet, err = templateSet.ParseGlob(filepath.Join(directory, "*.layout.tmpl"))

		if err != nil {
			return nil, err
		}

		templateSet, err = templateSet.ParseGlob(filepath.Join(directory, "*.partial.tmpl"))

		if err != nil {
			return nil, err
		}

		cache[name] = templateSet
	}

	return cache, nil
}
