package Hangman

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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

	// Utilisez le token d'accès pour faire une requête vers l'API Spotify
	artistName := "michael_jackson"
	artist, err := searchArtist(artistName, token)
	if err != nil {
		log.Fatalf("Erreur lors de la recherche de l'artiste: %v", err)
	}

	// Afficher les informations sur l'artiste
	fmt.Printf("Informations sur l'artiste:\nNom: %s\nType: %s\nPopularité: %d\nFollowers: %d\n",
		artist.Name, artist.Type, artist.Popularity, artist.Followers.Total)
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
		return "", fmt.Errorf("Token d'accès introuvable")
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
		return Artist{}, fmt.Errorf("Aucun artiste trouvé")
	}

	items, ok := artists["items"].([]interface{})
	if !ok || len(items) == 0 {
		return Artist{}, fmt.Errorf("Aucun artiste trouvé")
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

type Artist struct {
	Name       string
	Type       string
	Popularity int
	Followers  struct {
		Total int
	}
}
