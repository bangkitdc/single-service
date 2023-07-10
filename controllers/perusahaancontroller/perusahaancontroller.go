package perusahaancontroller

import (
	"api/helper"
	"api/models"
	"encoding/json"
	"net/http"
	"unicode"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetPerusahaans(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	queryParams := r.URL.Query()
	q := queryParams.Get("q")

	// Get Perusahaan
	var perusahaans []models.Perusahaan
	query := models.DB

	if q != "" {
		query = query.Where("nama = ? OR kode = ?", q, q)
	}

	err := query.Find(&perusahaans).Error
	if err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	// Find always returns false, so have to do it manually
	if len(perusahaans) == 0 {
		helper.ResponseJSON(w, http.StatusNotFound, "success", "No data available", []interface{}{})
		return
	}

	// Construct the response data
	var responseData []PerusahaanResponse
	for _, perusahaan := range perusahaans {
		item := PerusahaanResponse{
			ID:      perusahaan.ID,
			Nama:    perusahaan.Nama,
			Alamat:  perusahaan.Alamat,
			No_Telp: perusahaan.No_Telp,
			Kode:    perusahaan.Kode,
		}
		responseData = append(responseData, item)
	}

	helper.ResponseJSON(w, http.StatusOK, "success", "Perusahaan retrieved successfully", responseData)
}

func GetPerusahaanByID(w http.ResponseWriter, r *http.Request) {
	// Extract the ID parameter from the request URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Retrieve the Perusahaan from the database using the ID
	var perusahaan models.Perusahaan

	// Check if the Perusahaan is found
	if err := models.DB.Where("id = ?", id).First(&perusahaan).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helper.ResponseJSON(w, http.StatusUnauthorized, "error", "No data available", nil)
			return

		default:
			helper.ResponseJSON(w, http.StatusInternalServerError, "error", err.Error(), nil)
			return
		}
	}

	// Create the response data
	responseData := PerusahaanResponse{
		ID:      perusahaan.ID,
		Nama:    perusahaan.Nama,
		Alamat:  perusahaan.Alamat,
		No_Telp: perusahaan.No_Telp,
		Kode:    perusahaan.Kode,
	}

	// Return the response JSON
	helper.ResponseJSON(w, http.StatusOK, "success", "Perusahaan retrieved successfully", responseData)
}

func CreatePerusahaan(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var requestBody PerusahaanRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helper.ResponseJSON(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		return
	}

	// Validate the required fields
	if helper.IsEmpty(requestBody.Nama) || helper.IsEmpty(requestBody.Alamat) || helper.IsEmpty(requestBody.No_Telp) || helper.IsEmpty(requestBody.Kode) {
		helper.ResponseJSON(w, http.StatusBadRequest, "error", "Missing required fields", nil)
		return
	}

	// Create a new Perusahaan instance
	newPerusahaan := models.Perusahaan{
		Nama:    requestBody.Nama,
		Alamat:  requestBody.Alamat,
		No_Telp: requestBody.No_Telp,
		Kode:    requestBody.Kode,
	}

	// Check contraint
	message := CheckConstraint(&newPerusahaan)
	if !helper.IsEmpty(message) {
		helper.ResponseJSON(w, http.StatusBadRequest, "error", message, nil)
		return
	}

	// Save the new Perusahaan to the database
	err = models.DB.Create(&newPerusahaan).Error
	if err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	// Create the response data
	responseData := PerusahaanResponse{
		ID:      newPerusahaan.ID,
		Nama:    newPerusahaan.Nama,
		Alamat:  newPerusahaan.Alamat,
		No_Telp: newPerusahaan.No_Telp,
		Kode:    newPerusahaan.Kode,
	}

	// Return the response JSON
	helper.ResponseJSON(w, http.StatusCreated, "success", "Perusahaan created successfully", responseData)
}

func UpdatePerusahaan(w http.ResponseWriter, r *http.Request) {
	// Extract the ID parameter from the request URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Parse the request body
	var requestBody PerusahaanRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helper.ResponseJSON(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		return
	}

	// Validate the required fields
	if helper.IsEmpty(requestBody.Nama) || helper.IsEmpty(requestBody.Alamat) || helper.IsEmpty(requestBody.No_Telp) || helper.IsEmpty(requestBody.Kode) {
		helper.ResponseJSON(w, http.StatusBadRequest, "error", "Missing required fields", nil)
		return
	}

	// Retrieve the existing Perusahaan from the database using the ID
	var perusahaan models.Perusahaan
	if err := models.DB.Where("id = ?", id).First(&perusahaan).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helper.ResponseJSON(w, http.StatusNotFound, "error", "Perusahaan not found", nil)
			return
		default:
			helper.ResponseJSON(w, http.StatusInternalServerError, "error", "Failed to retrieve Perusahaan", nil)
			return
		}
	}

	// Update the fields of the existing Perusahaan with the new values
	perusahaan.Nama = requestBody.Nama
	perusahaan.Alamat = requestBody.Alamat
	perusahaan.No_Telp = requestBody.No_Telp
	perusahaan.Kode = requestBody.Kode

	// Check contraint
	message := CheckConstraint(&perusahaan)
	if !helper.IsEmpty(message) {
		helper.ResponseJSON(w, http.StatusBadRequest, "error", message, nil)
		return
	}

	// Save the updated Perusahaan back to the database
	if err := models.DB.Save(&perusahaan).Error; err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, "error", "Failed to update Perusahaan", nil)
		return
	}

	// Create the response data
	responseData := PerusahaanResponse{
		ID:      perusahaan.ID,
		Nama:    perusahaan.Nama,
		Alamat:  perusahaan.Alamat,
		No_Telp: perusahaan.No_Telp,
		Kode:    perusahaan.Kode,
	}

	// Return the response JSON
	helper.ResponseJSON(w, http.StatusOK, "success", "Perusahaan updated successfully", responseData)
}

func DeletePerusahaan(w http.ResponseWriter, r *http.Request) {
	// Extract the ID parameter from the request URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Retrieve the Perusahaan from the database using the ID
	var perusahaan models.Perusahaan

	// Check if the Perusahaan exists
	if err := models.DB.Where("id = ?", id).First(&perusahaan).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helper.ResponseJSON(w, http.StatusNotFound, "error", "perusahaan not found", nil)
			return

		default:
			helper.ResponseJSON(w, http.StatusInternalServerError, "error", err.Error(), nil)
			return
		}
	}

	// Delete the Perusahaan from the database
	if err := models.DB.Delete(&perusahaan).Error; err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	// Create the response data
	responseData := PerusahaanResponse{
		ID:      perusahaan.ID,
		Nama:    perusahaan.Nama,
		Alamat:  perusahaan.Alamat,
		No_Telp: perusahaan.No_Telp,
		Kode:    perusahaan.Kode,
	}

	helper.ResponseJSON(w, http.StatusOK, "success", "Perusahaan deleted successfully", responseData)
}

func CheckConstraint(responseBody *models.Perusahaan) string {
	if len(responseBody.Kode) != 3 {
		return "Kode (pajak) must be consists of 3 capital letters"
	}

	for _, char := range responseBody.Kode {
		if !unicode.IsUpper(char) {
			return "Kode (pajak) must be consists of 3 capital letters"
		}
	}

	return ""
}
