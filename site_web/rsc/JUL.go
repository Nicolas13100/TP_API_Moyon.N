package API

import (
	"encoding/base64"
	"fmt"
	"log"
)

func JUL() {
	clientID := "9b51a859f77e4bbda1729134d73e6676"
	clientSecret := "e22dafb4d6344f7d9704f034690f0a8c"

	// Encodez les informations d'identification du client pour l'authentification Basic << en français parceque pouquoi pas
	authHeader := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	// Obtenez un token d'accès OAuth2 << et lui aussi tien \o/
	token, err := getAccessToken(authHeader)
	if err != nil {
		log.Fatalf("Impossible d'obtenir un token: %v", err)
	}

	// Utilisez le token d'accès pour faire une requête vers l'API Spotify << sinon ba sa marche pas
	artistName := "JUL"

	artistID, err := searchArtistID(artistName, token)
	if err != nil {
		fmt.Println("Error searching for artist:", err)
		return
	}

	// Get the artist's albums << as asked
	albums, err := getArtistAlbums(artistID, token)
	if err != nil {
		fmt.Println("Error fetching albums:", err)
		return
	}

	// Display album information << cli version because, can be needed <<__>>
	for _, album := range albums {
		fmt.Println("Album Name:", album.Name)
		fmt.Println("Cover Image:", album.CoverImage)
		fmt.Println("Release Date:", album.ReleaseDate)
		fmt.Println("Number of Songs:", album.NumberOfSongs)
		fmt.Println("----------------------------")
	}
}

func getJULAlbums() ([]Album, error) {
	// use basic to make it cli too, just in case, for debug
	// found later i could have made all in one, but left like that since it's working :-D
	JUL()
	// Start of code
	clientID := "9b51a859f77e4bbda1729134d73e6676"
	clientSecret := "e22dafb4d6344f7d9704f034690f0a8c"

	// Encode the client's credentials for Basic authentication << asked by web documentation
	authHeader := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	// Get an OAuth2 access token << needed for query to work
	token, err := getAccessToken(authHeader)
	if err != nil {
		return nil, fmt.Errorf("unable to get token: %v", err)
	}

	// Search for the artist's ID << may be made interactive, not asked for now
	artistName := "JUL"
	artistID, err := searchArtistID(artistName, token)
	if err != nil {
		return nil, fmt.Errorf("error searching for artist: %v", err)
	}

	// Get the artist's albums ><><
	albums, err := getArtistAlbums(artistID, token)
	if err != nil {
		return nil, fmt.Errorf("error fetching albums: %v", err)
	}

	return albums, nil
}
