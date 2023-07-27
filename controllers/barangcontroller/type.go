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
	Nama          string    `json:"nama"`
	Harga         int       `json:"harga"`
	Stok          int       `json:"stok"`
	Kode          string    `json:"kode"`
	Perusahaan_ID uuid.UUID `json:"perusahaan_id"`
}

type PaginatedResponse struct {
	Data []BarangResponse `json:"data"`
	Meta Pagination       `json:"meta"`
}

type Pagination struct {
	Total        int64 `json:"total"`
	Current_Page int   `json:"current_page"`
	Total_Pages  int   `json:"total_pages"`
}

type BarangRequestQuantity struct {
	Quantity *int `json:"quantity"`
}
