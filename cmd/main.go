package main

import (
	"net/http"

	"urlshortener/handlers"
	"urlshortener/storage"

	"github.com/go-chi/chi"
)

func main() {

	storage.ConnectDB()
	route := chi.NewRouter()

	//implement routes

	route.Get("/urls/{alias}", handlers.RedirectUrl)
	route.Get("/urls", handlers.GetAllUrls)
	route.Post("/urls", handlers.MakeShortUrl)
	route.Delete("/urls/{alias}", handlers.DeleteUrl)
	route.Delete("/urls", handlers.DeleteAllUrls)

	//start the server
	http.ListenAndServe(":8080", route)

}
