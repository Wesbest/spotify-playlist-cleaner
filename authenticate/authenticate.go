package authenticate

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const redirectURI = "http://localhost:8080"

var (
	auth  *spotifyauth.Authenticator
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func init() {
	// Load the environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get Spotify credentials from the environment
	clientID := os.Getenv("SPOTIFY_ID")
	clientSecret := os.Getenv("SPOTIFY_SECRET")

	// Initialize the Spotify authenticator with client ID and secret
	auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithClientID(clientID), spotifyauth.WithClientSecret(clientSecret),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate))
}

func StartAuth() {
	// Handle favicon requests
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	// Handle the callback from Spotify
	http.HandleFunc("/callback", completeAuth)

	// Handle the root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
		fmt.Fprintln(w, "Welcome to the Spotify Authentication Service!")
	})

	// Start the HTTP server
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Print the authentication URL and wait for user to authenticate
	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// Wait for the authenticated client
	client := <-ch
	log.Println("Client received, fetching user info")

	// Retrieve and print user information
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	log.Println("Received callback request")

	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
		return
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
		return
	}

	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	log.Println("Callback completed successfully, sending client to channel")
	ch <- client
}
