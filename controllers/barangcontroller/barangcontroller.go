package barangcontroller

import (
	"api/helper"
	"api/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetBarangs(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	queryParams := r.URL.Query()
	q := queryParams.Get("q")
	perusahaan := queryParams.Get("perusahaan")

	// Get Barang
	var barangs []models.Barang
	query := models.DB

	if q != "" {
		query = query.Where("nama = ? OR kode = ?", q, q)
	}

	if perusahaan != "" {
		query = query.Where("perusahaan_id = ?", perusahaan)
	}

	err := query.Find(&barangs).Error
	if err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	// Find always returns false, so have to do it manually
	if len(barangs) == 0 {
		helper.ResponseJSON(w, http.StatusNotFound, "success", "No data available", []interface{}{})
		return
	}

	// Construct the response data
	var responseData []BarangResponse
	for _, barang := range barangs {
		item := BarangResponse{
			ID:            barang.ID,
			Harga:         barang.Harga,
			Nama:          barang.Nama,
			Stok:          barang.Stok,
			Kode:          barang.Kode,
			Perusahaan_ID: barang.Perusahaan_ID,
		}
		responseData = append(responseData, item)
	}

	helper.ResponseJSON(w, http.StatusOK, "success", "Barang retrieved successfully", responseData)
}

func GetBarangByID(w http.ResponseWriter, r *http.Request) {
	// Extract the ID parameter from the request URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Retrieve the Barang from the database using the ID
	var barang models.Barang

	// Check if the Barang is found
	if err := models.DB.Where("id = ?", id).First(&barang).Error; err != nil {
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
	responseData := BarangResponse{
		ID:            barang.ID,
		Harga:         barang.Harga,
		Nama:          barang.Nama,
		Stok:          barang.Stok,
		Kode:          barang.Kode,
		Perusahaan_ID: barang.Perusahaan_ID,
	}

	// Return the response JSON
	helper.ResponseJSON(w, http.StatusOK, "success", "Barang retrieved successfully", responseData)
}

func CreateBarang(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var requestBody BarangRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helper.ResponseJSON(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		return
	}

	// Validate the required fields
	if helper.IsEmpty(requestBody.Nama) || requestBody.Harga == nil || requestBody.Stok == nil || helper.IsEmpty(requestBody.Perusahaan_ID.String()) || helper.IsEmpty(requestBody.Kode) {
		helper.ResponseJSON(w, http.StatusBadRequest, "error", "Missing required fields", nil)
		return
	}

	// Create a new Barang instance
	newBarang := models.Barang{
		Nama:          requestBody.Nama,
		Harga:         *requestBody.Harga,
		Stok:          *requestBody.Stok,
		Perusahaan_ID: requestBody.Perusahaan_ID,
		Kode:          requestBody.Kode,
	}

	// Check contraint
	message := CheckConstraint(&newBarang, true)
	if !helper.IsEmpty(message) {
		helper.ResponseJSON(w, http.StatusBadRequest, "error", message, nil)
		return
	}

	// Save the new Barang to the database
	err = models.DB.Create(&newBarang).Error
	if err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	// Create the response data
	responseData := BarangResponse{
		ID:            newBarang.ID,
		Harga:         newBarang.Harga,
		Nama:          newBarang.Nama,
		Stok:          newBarang.Stok,
		Kode:          newBarang.Kode,
		Perusahaan_ID: newBarang.Perusahaan_ID,
	}

	// Return the response JSON
	helper.ResponseJSON(w, http.StatusCreated, "success", "Barang created successfully", responseData)
}

func UpdateBarang(w http.ResponseWriter, r *http.Request) {
	// Extract the ID parameter from the request URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Parse the request body
	var requestBody BarangRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helper.ResponseJSON(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		return
	}

	// Validate the required fields
	if helper.IsEmpty(requestBody.Nama) || requestBody.Harga == nil || requestBody.Stok == nil || helper.IsEmpty(requestBody.Perusahaan_ID.String()) || helper.IsEmpty(requestBody.Kode) {
		helper.ResponseJSON(w, http.StatusBadRequest, "error", "Missing required fields", nil)
		return
	}

	// Retrieve the existing Barang from the database using the ID
	var barang models.Barang
	if err := models.DB.Where("id = ?", id).First(&barang).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helper.ResponseJSON(w, http.StatusNotFound, "error", "Barang not found", nil)
			return
		default:
			helper.ResponseJSON(w, http.StatusInternalServerError, "error", "Failed to retrieve Barang", nil)
			return
		}
	}

	// Update the fields of the existing Barang with the new values
	barang.Nama = requestBody.Nama
	barang.Harga = *requestBody.Harga
	barang.Stok = *requestBody.Stok
	barang.Perusahaan_ID = requestBody.Perusahaan_ID
	barang.Kode = requestBody.Kode

	// Check contraint
	message := CheckConstraint(&barang, false)
	if !helper.IsEmpty(message) {
		helper.ResponseJSON(w, http.StatusBadRequest, "error", message, nil)
		return
	}

	// Save the updated Barang back to the database
	if err := models.DB.Save(&barang).Error; err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, "error", "Failed to update Barang", nil)
		return
	}

	// Create the response data
	responseData := BarangResponse{
		ID:            barang.ID,
		Nama:          barang.Nama,
		Harga:         barang.Harga,
		Stok:          barang.Stok,
		Kode:          barang.Kode,
		Perusahaan_ID: barang.Perusahaan_ID,
	}

	// Return the response JSON
	helper.ResponseJSON(w, http.StatusOK, "success", "Barang updated successfully", responseData)
}

func DeleteBarang(w http.ResponseWriter, r *http.Request) {
	// Extract the ID parameter from the request URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Retrieve the Barang from the database using the ID
	var barang models.Barang

	// Check if the Barang exists
	if err := models.DB.Where("id = ?", id).First(&barang).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helper.ResponseJSON(w, http.StatusNotFound, "error", "Barang not found", nil)
			return

		default:
			helper.ResponseJSON(w, http.StatusInternalServerError, "error", err.Error(), nil)
			return
		}
	}

	// Delete the Barang from the database
	if err := models.DB.Delete(&barang).Error; err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	// Create the response data
	responseData := BarangResponse{
		ID:            barang.ID,
		Harga:         barang.Harga,
		Nama:          barang.Nama,
		Stok:          barang.Stok,
		Kode:          barang.Kode,
		Perusahaan_ID: barang.Perusahaan_ID,
	}

	helper.ResponseJSON(w, http.StatusOK, "success", "Barang deleted successfully", responseData)
}

func CheckConstraint(responseBody *models.Barang, checkKode bool) string {
	var message string

	// Check Harga constraint
	if *&responseBody.Harga <= 0 {
		message += "Harga must be > 0"
	}

	// Check Stok constraint
	if *&responseBody.Stok < 0 {
		if !helper.IsEmpty(message) {
			message += ", "
		}
		message += "Stok must be >= 0"
	}

	// Check Perusahaan_ID validity
	var perusahaan models.Perusahaan
	if err := models.DB.Where("id = ?", responseBody.Perusahaan_ID).First(&perusahaan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if !helper.IsEmpty(message) {
				message += ", "
			}
			message += "Invalid perusahaan_id"
		} else {
			// Handle other database query errors
			log.Printf("Error querying Perusahaan: %v", err)
		}
	}

	// Check Kode uniqueness
	existingBarang := models.Barang{}
	if checkKode {
		if err := models.DB.Where("kode = ?", responseBody.Kode).First(&existingBarang).Error; err == nil {
			if !helper.IsEmpty(message) {
				message += ", "
			}
			message += "Kode must be unique"
		}
	} else {
		if err := models.DB.Where("kode = ? AND id != ?", responseBody.Kode, responseBody.ID).First(&existingBarang).Error; err == nil {
			if !helper.IsEmpty(message) {
				message += ", "
			}
			message += "Kode must be unique"
		}
	}

	return message
}
