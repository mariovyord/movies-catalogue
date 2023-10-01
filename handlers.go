package main

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type TemplateData struct {
	Title    string
	Content  template.HTML
	Username string
}

func handleMovieDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handleMovieDetailsPage(w, r)
	} else if r.Method == "DELETE" {
		handleMovieDelete(w, r)
	}
}

func handleAddMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		handleSubmitMovie(w, r)
	} else if r.Method == "GET" {
		handleAddMoviePage(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleIndexPage(w http.ResponseWriter, r *http.Request) {
	films := getMovies()
	renderPage(w, r)("Movie List", "index", films)
}

func handleMoviesPage(w http.ResponseWriter, r *http.Request) {
	films := getMovies()
	renderPage(w, r)("Movies", "movies", films)
}

func handleAddMoviePage(w http.ResponseWriter, r *http.Request) {
	renderPage(w, r)("Add Movie", "add-movie", nil)
}

func handleMovieDetailsPage(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/details/"):]

	movie, err := getMovieById(id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	renderPage(w, r)("Movie Details", "movie-details", movie)
}

func handleMovieDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/details/"):]

	err := deleteMovie(id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handleSignInPage(w, r)
	} else if r.Method == "POST" {
		handleSignInSubmit(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleSignInPage(w http.ResponseWriter, r *http.Request) {
	renderPage(w, r)("Sign In", "sign-in", nil)
}

func renderPage(w http.ResponseWriter, r *http.Request) func(title string, templateName string, templateDate interface{}) {
	return func(title string, templateName string, templateDate interface{}) {
		baseTemplate := template.Must(template.ParseFiles("./static/base.html"))
		tmpl := template.Must(template.ParseFiles("./static/" + templateName + ".html"))

		data := TemplateData{
			Title:    title,
			Username: GetUsernameFromCookie(r),
			Content:  "",
		}

		var buf bytes.Buffer
		err := tmpl.Execute(&buf, templateDate)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		data.Content = template.HTML(buf.String())

		err = baseTemplate.Execute(w, data)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func handleSubmitMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request to POST /add-movie/")

	title := r.PostFormValue("title")
	year, err := strconv.Atoi(r.PostFormValue("year"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	genre := r.PostFormValue("genre")
	director := r.PostFormValue("director")
	actors := strings.Split(r.PostFormValue("actors"), ",")

	film := Movie{Id: uuid.NewString(), Title: title, Director: director, Genre: genre, Year: year, Actors: actors}

	movie, err := addMovie(film)
	if err != nil {
		fmt.Println(err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("500 - Something bad happened!"))
		return
	}

	tmpl := template.Must(template.ParseFiles("./static/movie-added.html"))
	err = tmpl.Execute(w, movie.Id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
