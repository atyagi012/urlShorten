package main

import (
	"database/sql"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomCode(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func shorten(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	code := randomCode(6)
	_, err := db.Exec("INSERT INTO urls (code, original_url) VALUES ($1, $2)", code, req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"code": code})
}

func redirect(w http.ResponseWriter, r *http.Request) {
	var url string
	err := db.QueryRow("SELECT original_url FROM urls WHERE code = $1", r.PathValue("code")).Scan(&url)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func main() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/shorten", shorten)
	mux.HandleFunc("GET /api/{code}", redirect)
	http.ListenAndServe(":8080", mux)
}
