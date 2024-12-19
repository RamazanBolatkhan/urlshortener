package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"urlshortener/storage"

	"github.com/go-chi/chi"
	"golang.org/x/exp/rand"
)

func RedirectUrl(w http.ResponseWriter, r *http.Request) {

	alias := chi.URLParam(r, "alias")

	if alias == "" {
		// Handle missing alias in the URL path
		http.Error(w, "Alias is required", http.StatusBadRequest)
		return
	}

	var url storage.URL
	if err := storage.DB.First(&url, "alias = ?", alias).Error; err != nil {
		http.Error(w, "Could not find original url", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.OriginalUrl, http.StatusMovedPermanently)

}

func GetAllUrls(w http.ResponseWriter, r *http.Request) {
	var urls []storage.URL
	if err := storage.DB.Find(&urls).Error; err != nil {
		http.Error(w, "Failed to get all Urls", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urls)
}

// Deletes the URL based on the Alias that requested
func DeleteUrl(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "alias")
	if alias == "" {
		http.Error(w, "Alias is required", http.StatusBadRequest)
		return
	}

	if err := storage.DB.Where("alias = ?", alias).Delete(&storage.URL{}).Error; err != nil {
		http.Error(w, "Alias not found", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("URL Successfully deleted"))
}

// Deletes all records, all urls
func DeleteAllUrls(w http.ResponseWriter, r *http.Request) {
	//Clearing the urls table
	if err := storage.DB.Exec("DELETE FROM urls").Error; err != nil {
		http.Error(w, "Failed to delete all URLs", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All Urls successfully deleted"))
}

// Creates new short url with auto generated alias
func MakeShortUrl(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OriginalUrl string `json:"original_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request form", http.StatusBadRequest)
		return
	}
	alias := generateAlias()

	// Create a new URL entry
	newURL := storage.URL{
		OriginalUrl: req.OriginalUrl,
		Alias:       alias,
		CreatedTime: time.Now(),
	}

	// Append the new URL entry to the url table
	if err := storage.DB.Create(&newURL).Error; err != nil {
		http.Error(w, "Failed to create URL", http.StatusInternalServerError)
		return
	}

	// Respond with the new short URL
	resp := map[string]string{
		"short_url":    "http://localhost:8080/" + alias,
		"original_url": req.OriginalUrl,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func generateAlias() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6) // 6-character alias
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
