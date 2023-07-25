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

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Conect DB
	models.ConnectDatabase()

	// Seed the database
	models.Seeder()

	// Create router
	router := mux.NewRouter()

	// Enable CORS
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"https://ohl-fe.vercel.app"})

	// Login
	router.HandleFunc("/login", authcontroller.Login).Methods("POST")

	// Protected routes (Middleware)
	api := router.PathPrefix("/").Subrouter()
	api.Use(middleware.JWTMiddleware)

	// User
	api.HandleFunc("/self", selfcontroller.GetSelf).Methods("GET")

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

	// Doesn't need authorization, for monolith-service purpose only
	// Pagination Barang
	router.HandleFunc("/barang-paginate", barangcontroller.GetBarangsWithPagination).Methods("GET")

	// Barang
	router.HandleFunc("/barang-noauth/{id}", barangcontroller.GetBarangByID).Methods("GET")
	router.HandleFunc("/barang-noauth/{id}", barangcontroller.UpdateBarangWithQuantity).Methods("PUT")
	router.HandleFunc("/barang-noauth-recommendation", barangcontroller.GetBarangsRecommendation).Methods("GET")

	// Perusahaan
	router.HandleFunc("/perusahaan-noauth/{id}", perusahaancontroller.GetPerusahaanByID).Methods("GET")

	// Apply CORS middleware
	corsHandler := handlers.CORS(headers, methods, origins)(router)

	// Start server
	log.Fatal(http.ListenAndServe(":8000", corsHandler))
}
