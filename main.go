package main

import (
	"api/controllers/authcontroller"
	"api/controllers/barangcontroller"
	"api/controllers/perusahaancontroller"
	"api/controllers/selfcontroller"
	"api/middleware"
	"api/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func main() {
	// Conect DB
	models.ConnectDatabase()

	// Create router
	router := mux.NewRouter()

	// Login
	router.HandleFunc("/login", authcontroller.Login).Methods("POST")

	// Protected routes
	api := router.PathPrefix("/").Subrouter()
	api.Use(middleware.JWTMiddleware)

	// User
	api.HandleFunc("/", selfcontroller.GetSelf).Methods("GET")

	// Barang
	api.HandleFunc("/barang", barangcontroller.GetBarangs).Methods("GET")
	api.HandleFunc("/barang/{id}", barangcontroller.GetBarangByID).Methods("GET")
	api.HandleFunc("/barang", barangcontroller.CreateBarang).Methods("POST")
	api.HandleFunc("/barang/{id}", barangcontroller.UpdateBarang).Methods("PUT")
	api.HandleFunc("/barang/{id}", barangcontroller.DeleteBarang).Methods("DELETE")

	// Perusahaan
	api.HandleFunc("/perusahaan", perusahaancontroller.GetPerusahaans).Methods("GET")
	api.HandleFunc("/perusahaan/{id}", perusahaancontroller.GetPerusahaanByID).Methods("GET")
	api.HandleFunc("/perusahaan", perusahaancontroller.CreatePerusahaan).Methods("POST")
	api.HandleFunc("/perusahaan/{id}", perusahaancontroller.UpdatePerusahaan).Methods("PUT")
	api.HandleFunc("/perusahaan/{id}", perusahaancontroller.DeletePerusahaan).Methods("DELETE")

	// Start server
	log.Fatal(http.ListenAndServe(":8000", jsonContentTypeMiddleware(router)))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
