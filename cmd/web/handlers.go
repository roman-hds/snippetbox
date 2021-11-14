package main

import (
	"fmt"
	// "html/template"
	"net/http"
	"strconv"
)

// Change the signature of the home handler so it is defined as a method against *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w) // Use the notFound() helper
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, snippet := range s {
		fmt.Fprintf(w, "%v\n", snippet)
	}

	// Initialize a slice containing the paths to the two files. Note that the
	// home.page.tmpl file must be the *first* file in the slice.
	// files := []string{
	// 	"./ui/html/home.page.tmpl",
	// 	"./ui/html/base.layout.tmpl",
	// 	"./ui/html/footer.partial.tmpl",
	// }

	// // Use the template.ParseFiles() function to read the files and store the
	// // templates in a template set. Notice that we can pass the slice of file paths
	// // as a variadic parameter?
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	// Because the home handler function is now a method against application
	// 	// it can access its fields, including the error logger. We'll write the log
	// 	// message to this instead of the standard logger.
	// 	// app.errorLog.Println(err.Error())
	// 	app.serverError(w, err) // Use the serverError() helper.
	// 	return
	// }

	// // We then use the Execute() method on the template set to write the template content as the response body.
	// // The last parameter to Execute() represents any dynamic data that we want to pass in, which for now we'll leave as nil.
	// err = ts.Execute(w, nil)
	// if err != nil {
	// 	// Also update the code here to use the error logger from the application struct.
	// 	// app.errorLog.Println(err.Error())
	// 	// http.Error(w, "Internal Server Error", 500)
	// 	app.serverError(w, err) // Use the serverError() helper
	// }
}

// Change the signature of the showSnippet handler so it is defined as a method
// against *application.
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		// http.NotFound(w, r)
		app.notFound(w) // Use the notFound() helper.
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// Change the signature of the showSnippet handler so it is defined as a method
// against *application.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		// http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		app.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper.
		return
	}

	// Create some variables holding dummy data. We'll remove these later on during the build
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := "7"

	// Pass the data to the SnippetModel.Insert() method, receiving the ID of the new record back
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
