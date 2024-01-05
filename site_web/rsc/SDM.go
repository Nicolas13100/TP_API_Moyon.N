package API

import (
	"encoding/base64"
	"fmt"
	"log"
)

func SDM() {
	clientID := "9b51a859f77e4bbda1729134d73e6676"
	clientSecret := "e22dafb4d6344f7d9704f034690f0a8c"

	// Encodez les informations d'identification du client pour l'authentification Basic
	authHeader := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	// Obtenez un token d'accès OAuth2
	token, err := getAccessToken(authHeader)
	if err != nil {
		log.Fatalf("Impossible d'obtenir un token: %v", err)
	}

	trackInfo, err := searchTrack("SDM", "Bolide allemand", token)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Display track information
	fmt.Println("Title:", trackInfo.Title)
	fmt.Println("Album Cover:", trackInfo.AlbumCover)
	fmt.Println("Album Name:", trackInfo.AlbumName)
	fmt.Println("Artist Name:", trackInfo.ArtistName)
	fmt.Println("Release Date:", trackInfo.ReleaseDate)
	fmt.Println("Spotify Link:", trackInfo.SpotifyLink)
}
