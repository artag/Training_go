package main

import (
	"html/template" // New import
	"path/filepath" // New import
	"snippetbox/internal/models"
)

// The holding structure for any dynamic data
// that we want to pass to HTML templates.
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

// Create cache fpr templates.
func newTemplateCache() (map[string]*template.Template, error) {
	// Init a new map for cache
	cache := map[string]*template.Template{}

	// Get a slice of all filepaths that match the pattern
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract file name (like 'home.tmpl') from the path
		name := filepath.Base(page)

		// Parse the base template file into a template set.
		ts, err := template.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// Parse (and add) a slice of filepaths for all partials into a template set.
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// Parse (and add) the page into a template set.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map.
		cache[name] = ts
	}

	return cache, nil
}
