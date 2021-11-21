package main

import "roman-hds/snippetbox/pkg/models"

// templateData acts as a holding struct for any dynamic data that we want to pass to HTML templates
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
