package main

import (
	API "API/site_web/rsc"
	"encoding/base64"
	"encoding/json"
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

func main() {
	clientID := "9b51a859f77e4bbda1729134d73e6676"
	clientSecret := "e22dafb4d6344f7d9704f034690f0a8c"

	// Encodez les informations d'identification du client pour l'authentification Basic
	authHeader := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	// Obtenez un token d'accès OAuth2
	token, err := getAccessToken(authHeader)
	if err != nil {
		log.Fatalf("Impossible d'obtenir un token: %v", err)
	}

}

func JUL() {
	clientID := "9b51a859f77e4bbda1729134d73e6676"
	clientSecret := "e22dafb4d6344f7d9704f034690f0a8c"

	// Encodez les informations d'identification du client pour l'authentification Basic
	authHeader := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	// Obtenez un token d'accès OAuth2
	token, err := getAccessToken(authHeader)
	if err != nil {
		log.Fatalf("Impossible d'obtenir un token: %v", err)
	}

	// Utilisez le token d'accès pour faire une requête vers l'API Spotify
	artistName := "JUL"
	artist, err := searchArtist(artistName, token)
	if err != nil {
		log.Fatalf("Erreur lors de la recherche de l'artiste: %v", err)
	}

	// Afficher les informations sur l'artiste
	fmt.Printf("Informations sur l'artiste:\nNom: %s\nType: %s\nPopularité: %d\nFollowers: %d\n",
		artist.Name, artist.Type, artist.Popularity, artist.Followers.Total)
	fmt.Println("----------------------------")

	artistID, err := searchArtistID(artistName, token)
	if err != nil {
		fmt.Println("Error searching for artist:", err)
		return
	}

	// Get the artist's albums
	albums, err := getArtistAlbums(artistID, token)
	if err != nil {
		fmt.Println("Error fetching albums:", err)
		return
	}

	// Display album information
	for _, album := range albums {
		fmt.Println("Album Name:", album.Name)
		fmt.Println("Cover Image:", album.CoverImage)
		fmt.Println("Release Date:", album.ReleaseDate)
		fmt.Println("Number of Songs:", album.NumberOfSongs)
		fmt.Println("----------------------------")
	}
}

func getAccessToken(authHeader string) (string, error) {
	tokenURL := "https://accounts.spotify.com/api/token"
	payload := "?grant_type=client_credentials"
	body := tokenURL + payload
	req, err := http.NewRequest("POST", body, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Basic "+authHeader)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	accessToken, ok := tokenResp["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("token d'accès introuvable")
	}

	return accessToken, nil
}

func searchArtist(artistName string, accessToken string) (API.Artist, error) {
	searchURL := "https://api.spotify.com/v1/search?type=artist&q=" + artistName

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return API.Artist{}, err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return API.Artist{}, err
	}
	defer resp.Body.Close()

	var searchResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return API.Artist{}, err
	}

	artists, ok := searchResp["artists"].(map[string]interface{})
	if !ok {
		return API.Artist{}, fmt.Errorf("aucun artiste trouvé")
	}

	items, ok := artists["items"].([]interface{})
	if !ok || len(items) == 0 {
		return API.Artist{}, fmt.Errorf("aucun artiste trouvé")
	}

	// Récupérer les détails du premier artiste trouvé
	artistData := items[0].(map[string]interface{})

	var artist API.Artist
	artist.Name = artistData["name"].(string)
	artist.Type = artistData["type"].(string)
	artist.Popularity = int(artistData["popularity"].(float64))

	followers := artistData["followers"].(map[string]interface{})
	artist.Followers.Total = int(followers["total"].(float64))

	return artist, nil
}

func getArtistAlbums(artistID string, accessToken string) ([]API.Album, error) {
	albumsURL := fmt.Sprintf("https://api.spotify.com/v1/artists/%s/albums", artistID)

	req, err := http.NewRequest("GET", albumsURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var albumsResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&albumsResp); err != nil {
		return nil, err
	}

	albumsData, ok := albumsResp["items"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("aucun album trouvé")
	}

	var albums []API.Album
	for _, album := range albumsData {
		albumData := album.(map[string]interface{})
		var newAlbum API.Album
		newAlbum.Name = albumData["name"].(string)
		images := albumData["images"].([]interface{})
		if len(images) > 0 {
			newAlbum.CoverImage = images[0].(map[string]interface{})["url"].(string)
		}
		newAlbum.ReleaseDate = albumData["release_date"].(string)
		tracks := albumData["total_tracks"].(float64)
		newAlbum.NumberOfSongs = int(tracks)
		albums = append(albums, newAlbum)
	}

	return albums, nil
}

func searchArtistID(artistName string, accessToken string) (string, error) {
	searchURL := "https://api.spotify.com/v1/search?type=artist&q=" + artistName

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var searchResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return "", err
	}

	artists, ok := searchResp["artists"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("aucun artiste trouvé")
	}

	items, ok := artists["items"].([]interface{})
	if !ok || len(items) == 0 {
		return "", fmt.Errorf("aucun artiste trouvé")
	}

	// Return the ID of the first artist found
	artistData := items[0].(map[string]interface{})
	return artistData["id"].(string), nil
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
