package API

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
)

func RUN() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/album/jul", julHandler)
	http.HandleFunc("/track/sdm", sdmHandler)
	http.HandleFunc("/gestion/jul", GjulHandler)
	http.HandleFunc("/gestion/sdm", GsdmHandler)

	// Serve static files from the "static" directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Print statement indicating server is running
	fmt.Println("Server is running on :8080 http://localhost:8080")

	// Start the server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}

func julHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "jul", nil)
}

func sdmHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "bolideAllemand", nil)
}

func GjulHandler(w http.ResponseWriter, r *http.Request) {}

func GsdmHandler(w http.ResponseWriter, r *http.Request) {}

func renderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	tmpl, err := template.New(tmplName).Funcs(template.FuncMap{"join": join}).ParseFiles("Template/" + tmplName + ".html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func join(s []string, sep string) string {
	return strings.Join(s, sep)
}
