// seeder.go

package models

import (
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
)

// Seeder inserts seed data into the database.
func Seeder() {
	// Seed Users
	seedUsers()

	// Seed Perusahaan
	seedPerusahaans()

	// Seed Barang
	seedBarangs()
}

func seedUsers() {
	// Static user data
	user := User{
		Username: "admin",
		Name:     "admin",
		Password: "admin",
	}

	if err := DB.Create(&user).Error; err != nil {
		log.Fatalf("Failed to seed user: %v", err)
	}
}

func seedPerusahaans() {
	// Generate and seed fake perusahaan data
	perusahaans := make([]Perusahaan, 0, 10)
	// Create a set to store used company names
	usedNames := make(map[string]bool)

	for len(perusahaans) < 10 {
		// Generate a unique company name for the perusahaan
		nama := strings.Title(gofakeit.Company())
		if !usedNames[nama] {
			usedNames[nama] = true
			perusahaan := Perusahaan{
				Nama:    "PT " + nama,
				Alamat:  "Jalan " + gofakeit.StreetName(),
				Kode:    generateRandomAlphabet(3),
				No_Telp: "08" + generateRandomDigits(10),
			}
			perusahaans = append(perusahaans, perusahaan)
		}
	}

	// Batch insert the generated perusahaans
	if err := DB.Create(&perusahaans).Error; err != nil {
		log.Fatalf("Failed to seed perusahaans: %v", err)
	}
}

func seedBarangs() {
	// Generate and seed fake barang data
	barangs := make([]Barang, 0, 20)
	// Fetch all Perusahaan IDs from the database
	var perusahaanIDs []uuid.UUID
	if err := DB.Model(&Perusahaan{}).Pluck("id", &perusahaanIDs).Error; err != nil {
		log.Fatalf("Failed to fetch Perusahaan IDs: %v", err)
	}

	// Create a map to store used names
	usedNames := make(map[string]bool)

	for len(barangs) < 20 {
		// Generate a unique name for the barang
		nama := strings.Title(gofakeit.NounCountable())
		if !usedNames[nama] {
			usedNames[nama] = true
			barang := Barang{
				Nama:          nama,
				Harga:         rand.Intn(100) * 1000,
				Stok:          rand.Intn(1000),
				Kode:          string(nama[0]) + generateRandomDigits(3),
				Perusahaan_ID: perusahaanIDs[rand.Intn(len(perusahaanIDs))],
			}
			barangs = append(barangs, barang)
		}
	}

	// Batch insert the generated barangs
	if err := DB.Create(&barangs).Error; err != nil {
		log.Fatalf("Failed to seed barang: %v", err)
	}
}

func generateRandomDigits(length int) string {
	var result string
	for i := 0; i < length; i++ {
		result += strconv.Itoa(rand.Intn(10))
	}
	return result
}

func generateRandomAlphabet(n int) string {
	var result string
	letterRunes := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for i := 0; i < n; i++ {
		result += string(letterRunes[rand.Intn(len(letterRunes))])
	}
	return result
}
