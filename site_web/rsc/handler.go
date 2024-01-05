package API

import (
	"fmt"
	"log"
	"net/http"
)

func RUN() {

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/album/jul", julHandler)
	http.HandleFunc("/track/sdm", sdmHandler)

	// Serve static files from the "site_web/static" directory << copied from hangman
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("site_web/static"))))

	// Print statement indicating server is running << same
	fmt.Println("Server is running on :8080 http://localhost:8080")

	// Start the server << same
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}

func julHandler(w http.ResponseWriter, r *http.Request) {
	// Call the function to retrieve album information and send them to the page <>.<>
	albums, err := getJULAlbums()
	if err != nil {
		http.Error(w, "Error retrieving albums", http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "jul", albums)
}

func sdmHandler(w http.ResponseWriter, r *http.Request) {
	// Call the function to retrieve track information and send them to the page <>.<>
	trackInfo, err := SDM()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "sdm", trackInfo)
}
