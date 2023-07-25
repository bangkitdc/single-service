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
	perusahaans := make([]Perusahaan, 10)
	for _, perusahaan := range perusahaans {
		nama := gofakeit.Company()
		perusahaan.Nama = "PT " + strings.Title(nama)
		perusahaan.Alamat = "Jalan " + gofakeit.StreetName()
		perusahaan.Kode = generateRandomAlphabet(3)

		// Manually generate the phone number with "08" prefix and 10 random digits
		perusahaan.No_Telp = "08" + generateRandomDigits(10)

		if err := DB.Create(&perusahaan).Error; err != nil {
			log.Fatalf("Failed to seed perusahaan: %v", err)
		}
	}
}

func seedBarangs() {
	// Generate and seed fake barang data
	barangs := make([]Barang, 20)
	// Fetch all Perusahaan IDs from the database
	var perusahaanIDs []uuid.UUID
	if err := DB.Model(&Perusahaan{}).Pluck("id", &perusahaanIDs).Error; err != nil {
		log.Fatalf("Failed to fetch Perusahaan IDs: %v", err)
	}

	for _, barang := range barangs {
		nama := gofakeit.NounCountable()
		barang.Nama = strings.Title(nama)
		barang.Harga = rand.Intn(100) * 1000
		barang.Stok = rand.Intn(1000)
		barang.Kode = string(barang.Nama[0]) + generateRandomDigits(3)
		barang.Perusahaan_ID = perusahaanIDs[rand.Intn(len(perusahaanIDs))]

		if err := DB.Create(&barang).Error; err != nil {
			log.Fatalf("Failed to seed barang: %v", err)
		}
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
