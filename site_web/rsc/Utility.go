package API

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

func searchArtist(artistName string, accessToken string) (Artist, error) {
	searchURL := "https://api.spotify.com/v1/search?type=artist&q=" + artistName

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return Artist{}, err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Artist{}, err
	}
	defer resp.Body.Close()

	var searchResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return Artist{}, err
	}

	artists, ok := searchResp["artists"].(map[string]interface{})
	if !ok {
		return Artist{}, fmt.Errorf("aucun artiste trouvé")
	}

	items, ok := artists["items"].([]interface{})
	if !ok || len(items) == 0 {
		return Artist{}, fmt.Errorf("aucun artiste trouvé")
	}

	// Récupérer les détails du premier artiste trouvé
	artistData := items[0].(map[string]interface{})

	var artist Artist
	artist.Name = artistData["name"].(string)
	artist.Type = artistData["type"].(string)
	artist.Popularity = int(artistData["popularity"].(float64))

	followers := artistData["followers"].(map[string]interface{})
	artist.Followers.Total = int(followers["total"].(float64))

	return artist, nil
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
