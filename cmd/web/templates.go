package main

import (
	"html/template"
	"path/filepath"
	"time"

	"zapmal/snippetbox/pkg/forms"
	"zapmal/snippetbox/pkg/models"
)

type TemplateData struct {
	CurrentYear int
	Form        *forms.Form
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

func humanDate(time time.Time) string {
	return time.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(directory string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(directory, "*.page.tmpl"))

	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		templateSet, err := template.New(name).Funcs(functions).ParseFiles(page)

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
