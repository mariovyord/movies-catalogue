package main

import "fmt"

var MOVIES_URL = "database/movies.json"

func getAllMovies() ([]Movie, error) {
	movies, err := get(MOVIES_URL)
	return movies, err
}

func getMovieById(id string) (*Movie, error) {
	movies, err := get(MOVIES_URL)
	if err != nil {
		return nil, err
	}

	for _, movie := range movies {
		if movie.Id == id {
			return &movie, nil
		}
	}

	return nil, nil
}

func addMovie(movie Movie) (Movie, error) {
	_, err := post(MOVIES_URL, movie)
	if err != nil {
		return Movie{}, nil
	}
	return movie, nil
}

func updateMovie(id string, movie Movie) error {
	_, err := patch(MOVIES_URL, movie)
	if err != nil {
		return err
	}

	return nil
}

func deleteMovie(id string) error {
	_, err := del(MOVIES_URL, id)
	if err != nil {
		return err
	}

	return nil
}

func getMovies() map[string][]Movie {
	allMovies, err := getAllMovies()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	films := map[string][]Movie{
		"Films": allMovies,
	}

	return films
}
