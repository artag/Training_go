package main

import "snippetbox/internal/models"

// The holding structure for any dynamic data
// that we want to pass to HTML templates.
type templateData struct {
	Snippet *models.Snippet
}
