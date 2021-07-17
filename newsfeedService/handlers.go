package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type categoryResponse struct {
	Categories []Category `json:"categories"`
}

type newsResponse struct {
	News []News `json:"news"`
}

func (srv service) categoryHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := srv.db.getCategories()
	if err != nil {
		log.Println(err)
		JSON(w, 500, "error fetching categories")
	}
	result := categoryResponse{
		Categories: append(make([]Category, 0), categories...),
	}
	JSON(w, 200, result)
}

func (srv service) newsHandler(w http.ResponseWriter, r *http.Request) {
	news, err := srv.db.getNews()
	if err != nil {
		log.Println(err)
		JSON(w, 500, "error fetching news")
	}
	result := newsResponse{
		News: append(make([]News, 0), news...),
	}
	JSON(w, 200, result)
}

func (srv service) newsByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	news, err := srv.db.getNewsByCategory(vars["category"])
	if err != nil {
		log.Println(err)
		JSON(w, 500, "error fetching news")
	}
	result := newsResponse{
		News: append(make([]News, 0), news...),
	}
	JSON(w, 200, result)
}
