package barangcontroller

import (
	"api/helper"
	"api/models"
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetBarangs(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	queryParams := r.URL.Query()
	q := queryParams.Get("q")
	perusahaanID := queryParams.Get("perusahaan")

	// Get Barang
	var barangs []models.Barang
	query := models.DB

	if q != "" {
		query = query.Where("nama ILIKE ? OR kode ILIKE ?", "%"+q+"%", "%"+q+"%")
	}

	if perusahaanID != "" {
		query = query.Where("perusahaan_id = ?", perusahaanID)
	}

	// Order by UUID
	query = query.Order("nama")

	err := query.Find(&barangs).Error
	if err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	// Find always returns false, so have to do it manually
	if len(barangs) == 0 {
		helper.ResponseJSON(w, http.StatusOK, "success", "No data available", []interface{}{})
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
			helper.ResponseJSON(w, http.StatusNotFound, "error", "No data available", nil)
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

// Constraint Checker
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
	existingBarangs := make([]models.Barang, 0)
	if checkKode {
		models.DB.Where("kode = ?", responseBody.Kode).Find(&existingBarangs)
	} else {
		models.DB.Where("kode = ? AND id != ?", responseBody.Kode, responseBody.ID).Find(&existingBarangs)
	}

	if len(existingBarangs) > 0 {
		// Records were found, so Kode is not unique
		if !helper.IsEmpty(message) {
			message += ", "
		}
		message += "Kode must be unique"
	}

	return message
}

// Pagination
func GetBarangsWithPagination(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	queryParams := r.URL.Query()
	q := queryParams.Get("q")
	perusahaanID := queryParams.Get("perusahaan")
	pageStr := queryParams.Get("page")

	// Convert page query parameter to an integer
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1 // Default page number if not provided or invalid
	}

	// Set the number of items per page
	limit := 8

	// Calculate the offset
	offset := (page - 1) * limit

	// Get Barang with pagination
	var barangs []models.Barang
	query := models.DB

	if q != "" {
		query = query.Where("nama ILIKE ? OR kode ILIKE ?", "%"+q+"%", "%"+q+"%")
	}

	if perusahaanID != "" {
		query = query.Where("perusahaan_id = ?", perusahaanID)
	}

	// Order by UUID
	query = query.Order("nama")

	// Get total count without applying pagination
	var totalCount int64
	query.Model(&models.Barang{}).Count(&totalCount)

	// Calculate the total number of pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	// Apply pagination
	query = query.Offset(offset).Limit(limit)

	err = query.Find(&barangs).Error
	if err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, "error", err.Error(), nil)
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

	// Construct the pagination metadata
	pagination := Pagination{
		Total:        totalCount,
		Current_Page: page,
		Total_Pages:  totalPages,
	}

	helper.ResponseJSON(w, http.StatusOK, "success", "Barang retrieved successfully", PaginatedResponse{
		Data: responseData,
		Meta: pagination,
	})
}

// Recommendation
func GetBarangsRecommendation(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	queryParams := r.URL.Query()
	nama := queryParams.Get("nama")
	except := queryParams.Get("except")

	// Find the barang with the given nama
	var barang models.Barang
	if err := models.DB.Where("nama LIKE ?", "%"+nama+"%").First(&barang).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helper.ResponseJSON(w, http.StatusNotFound, "error", "Barang not found", nil)
			return
		default:
			helper.ResponseJSON(w, http.StatusInternalServerError, "error", "Failed to retrieve Barang", nil)
			return
		}
	}

	// Get Barangs with the same perusahaan_id (up to a maximum of 4 random items)
	var barangs []models.Barang
	err2 := models.DB.Where("perusahaan_id = ? AND nama != ?", barang.Perusahaan_ID, except).Order("RANDOM()").Limit(4).Find(&barangs).Error
	if err2 != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, "error", err2.Error(), nil)
		return
	}

	// Find always returns false, so have to do it manually
	if len(barangs) == 0 {
		helper.ResponseJSON(w, http.StatusOK, "success", "No data available", []interface{}{})
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

// Update Stok
func UpdateBarangWithQuantity(w http.ResponseWriter, r *http.Request) {
	// Extract the ID parameter from the request URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Parse the request body
	var requestBody BarangRequestQuantity

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helper.ResponseJSON(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		return
	}

	// Validate the required fields
	if requestBody.Quantity == nil {
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
	barang.Stok -= *requestBody.Quantity

	// Check contraint
	if barang.Stok < 0 {
		helper.ResponseJSON(w, http.StatusBadRequest, "error", "Check the stock again", nil)
		return
	}

	// Save the updated Barang back to the database
	if err := models.DB.Save(&barang).Error; err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, "error", "Failed to Checkout Barang", nil)
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
