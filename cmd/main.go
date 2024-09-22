package main

import (
	"log"
	"spotify-playlist-cleaner/authenticate"
	"spotify-playlist-cleaner/playlist"
)

func main() {
	// Declare the playlist name as a variable
	playlistName := "About last month"

	// Start the Spotify authentication process
	client := authenticate.StartAuth() // Now it returns a valid *spotify.Client

	// After authentication, update the playlist
	err := playlist.UpdatePlaylist(client, playlistName)
	if err != nil {
		log.Fatalf("Failed to update playlist: %v", err)
	}

}
