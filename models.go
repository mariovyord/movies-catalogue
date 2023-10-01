package main

type Movie struct {
	Id       string   `json:"id"`
	Title    string   `json:"title"`
	Year     int      `json:"year"`
	Genre    string   `json:"genre"`
	Director string   `json:"director"`
	Actors   []string `json:"actors"`
}
