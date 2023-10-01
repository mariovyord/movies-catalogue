package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/add-movie/", handleAddMovie)
	mux.HandleFunc("/details/", handleMovieDetails)
	mux.HandleFunc("/movies/", handleMoviesPage)
	mux.HandleFunc("/sign-in/", handleSignIn)
	mux.HandleFunc("/sign-out/", handleSignOut)
	mux.HandleFunc("/", handleIndexPage)

	fmt.Printf("Server running at http://localhost:3001\n")
	log.Fatal(http.ListenAndServe(":3001", handleLogging(mux)))
}
