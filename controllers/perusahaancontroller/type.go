package perusahaancontroller

import uuid "github.com/satori/go.uuid"

type PerusahaanRequest struct {
	Nama    string `json:"nama"`
	Alamat  string `json:"alamat"`
	No_Telp string `json:"no_telp"`
	Kode    string `json:"kode"`
}

type PerusahaanResponse struct {
	ID      uuid.UUID `json:"id"`
	Nama    string    `json:"nama"`
	Alamat  string    `json:"alamat"`
	No_Telp string    `json:"no_telp"`
	Kode    string    `json:"kode"`
}
