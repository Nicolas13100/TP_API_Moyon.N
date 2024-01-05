package API

import (
	"encoding/base64"
	"fmt"
	"log"
)

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
