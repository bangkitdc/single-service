package barangcontroller

import uuid "github.com/satori/go.uuid"

type BarangRequest struct {
	Nama          string    `json:"nama"`
	Harga         *int      `json:"harga"`
	Stok          *int      `json:"stok"`
	Perusahaan_ID uuid.UUID `json:"perusahaan_id"`
	Kode          string    `json:"kode"`
}

type BarangResponse struct {
	ID            uuid.UUID `json:"id"`
	Harga         int       `json:"harga"`
	Nama          string    `json:"nama"`
	Stok          int       `json:"stok"`
	Kode          string    `json:"kode"`
	Perusahaan_ID uuid.UUID `json:"perusahaan_id"`
}
