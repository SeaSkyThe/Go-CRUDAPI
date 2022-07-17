package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

// Function that READ all the movies
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //With the response writer, we set the content type to json
	json.NewEncoder(w).Encode(movies)                  //Then, send all the movies as a json object
}

// Function that READ a specific movie
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //With the response writer, we set the content type to json
	params := mux.Vars(r)                              //Getting the params
	for _, movie := range movies {                     //For each movie in "database"
		if movie.ID == params["id"] { //If we found a movie that has the same ID, that the one sent as parameter, we READ
			json.NewEncoder(w).Encode(movie)
			return
		}
	}

}

// Function that DELETE a specific movie
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //With the response writer, we set the content type to json
	params := mux.Vars(r)                              //Getting the params
	for index, movie := range movies {                 //For each movie in "database"
		if movie.ID == params["id"] { //If we found a movie that has the same ID, that the one sent as parameter, we delete
			movies = append(movies[:index], movies[index+1:]...) //Delete
			break
		}
	}
	json.NewEncoder(w).Encode(movies)

}

// Function that CREATES a movie
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //With the response writer, we set the content type to json
	var movie Movie                                    //Create the movie variable
	_ = json.NewDecoder(r.Body).Decode(&movie)         //Get the infos about the movie sent in the body, and put into the movie variable
	movie.ID = strconv.Itoa(rand.Intn(1000000000))     //Define a unique ID
	movies = append(movies, movie)                     //Append the new movie to the "database"
	json.NewEncoder(w).Encode(movie)
}

// Function that UPDATE a movie
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //Set the content type to json
	params := mux.Vars(r)                              //Get params
	for index, movie := range movies {                 //Search for the movie in "database"
		if movie.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...) //Delete de movie with de ID sent
			var movie Movie                                      //Add a new movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}

}
func main() {
	r := mux.NewRouter() //Router

	movies = append(movies, Movie{ID: "1", Isbn: "438227", Title: "Movie One", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "454555", Title: "Movie Two", Director: &Director{Firstname: "Steve", Lastname: "Smith"}})
	//Functions (routes) that the router handles
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8000...")
	log.Fatal(http.ListenAndServe(":8000", r))

}
