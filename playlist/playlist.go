package playlist

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zmb3/spotify/v2"
)

// FindPlaylistByName finds a playlist by its name
func FindPlaylistByName(client *spotify.Client, name string) (*spotify.SimplePlaylist, error) {
	// Get the current user's playlists
	playlists, err := client.CurrentUsersPlaylists(context.Background())
	if err != nil {
		return nil, err
	}

	// Loop through playlists to find the one with the specified name
	for _, playlist := range playlists.Playlists {
		if playlist.Name == name {
			return &playlist, nil
		}
	}

	return nil, fmt.Errorf("playlist '%s' not found", name)
}

// GetPlaylistTracks retrieves all tracks from a specific playlist
func GetPlaylistTracks(client *spotify.Client, playlistID spotify.ID) ([]spotify.PlaylistTrack, error) {
	// Get the playlist's tracks
	tracks, err := client.GetPlaylistTracks(context.Background(), playlistID)
	if err != nil {
		return nil, err
	}

	return tracks.Tracks, nil
}

// FilterOldTracks filters out tracks that are older than one month
func FilterOldTracks(tracks []spotify.PlaylistTrack) []spotify.PlaylistTrack {
	var recentTracks []spotify.PlaylistTrack
	oneMonthAgo := time.Now().AddDate(0, -1, 0) // One month ago

	for _, track := range tracks {
		// Parse the AddedAt field, which is a string, into a time.Time object
		addedAt, err := time.Parse(time.RFC3339, track.AddedAt)
		if err != nil {
			log.Printf("Error parsing time for track %s: %v", track.Track.Name, err)
			continue
		}

		// Check if the track was added after one month ago
		if addedAt.After(oneMonthAgo) {
			recentTracks = append(recentTracks, track)
		}
	}

	return recentTracks
}

// RemoveOldTracks removes the tracks older than one month from the playlist
func RemoveOldTracks(client *spotify.Client, playlistID spotify.ID, tracks []spotify.PlaylistTrack) error {
	for _, track := range tracks {
		_, err := client.RemoveTracksFromPlaylist(context.Background(), playlistID, track.Track.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdatePlaylist updates the specified playlist by removing tracks older than one month
func UpdatePlaylist(client *spotify.Client, playlistName string) error {
	// Step 1: Find the playlist by name
	playlist, err := FindPlaylistByName(client, playlistName)
	if err != nil {
		return err
	}

	// Step 2: Get the playlist's tracks
	tracks, err := GetPlaylistTracks(client, playlist.ID)
	if err != nil {
		return err
	}

	// Step 3: Filter tracks older than one month
	oldTracks := FilterOldTracks(tracks)

	// Step 4: Remove old tracks from the playlist
	err = RemoveOldTracks(client, playlist.ID, oldTracks)
	if err != nil {
		return err
	}

	fmt.Println("Successfully updated playlist by removing old tracks.")
	return nil
}
