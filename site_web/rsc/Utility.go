package API

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

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

func getArtistAlbums(artistID string, accessToken string) ([]Album, error) {
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

	var albums []Album
	for _, album := range albumsData {
		albumData := album.(map[string]interface{})
		var newAlbum Album
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

func searchTrack(artistName, trackName, accessToken string) (*TrackInfo, error) {
	// Construct the search query
	query := fmt.Sprintf("track:%s artist:%s", trackName, artistName)
	query = strings.ReplaceAll(query, " ", "%20")

	searchURL := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track", query)

	req, err := http.NewRequest("GET", searchURL, nil)
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

	var searchResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	tracks, ok := searchResp["tracks"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("track not found")
	}

	items, ok := tracks["items"].([]interface{})
	if !ok || len(items) == 0 {
		return nil, fmt.Errorf("track not found")
	}

	// Extract the first track found
	trackData := items[0].(map[string]interface{})

	// Extract relevant track information
	trackInfo := &TrackInfo{}
	trackInfo.Title = trackData["name"].(string)
	trackInfo.AlbumCover = trackData["album"].(map[string]interface{})["images"].([]interface{})[0].(map[string]interface{})["url"].(string)
	trackInfo.AlbumName = trackData["album"].(map[string]interface{})["name"].(string)
	artists := trackData["artists"].([]interface{})
	if len(artists) > 0 {
		trackInfo.ArtistName = artists[0].(map[string]interface{})["name"].(string)
	}
	trackInfo.ReleaseDate = trackData["album"].(map[string]interface{})["release_date"].(string)
	trackInfo.SpotifyLink = trackData["external_urls"].(map[string]interface{})["spotify"].(string)

	return trackInfo, nil
}

func renderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {

	tmpl, err := template.New(tmplName).Funcs(template.FuncMap{"join": join}).ParseFiles("site_web/Template/" + tmplName + ".html")
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
