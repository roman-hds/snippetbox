package main

import (
	"html/template"
	"path/filepath"
	"roman-hds/snippetbox/pkg/models"
	"time"
)

// templateData acts as a holding struct for any dynamic data that we want to pass to HTML templates
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() to get a slice of all felipaths with the extention '.page.tmpl'.
	// This essentially gives us a slice of all the 'page' templates for the application
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Loop trough the pages one-by-one
	for _, page := range pages {
		// Extract the file names (like 'home.page.tmpl') from the full file path
		// and assign it to the name variable
		name := filepath.Base(page)

		// Parse the page template file into a template set
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// Use the ParseGlob() to add any 'layout' templates to the template set
		// in our case, it's just the 'base' layout' at the moment
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'partial' templates to the
		// template set (in our case, it's just the 'footer' partial at the moment).
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache, using the name of the page
		// (like 'home.page.tmpl') as the key
		cache[name] = ts
	}

	// Return the map
	return cache, nil
}
