package models

import uuid "github.com/satori/go.uuid"

type Barang struct {
	ID            uuid.UUID  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Nama          string     `gorm:"varchar(255);not null" json:"nama"`
	Harga         int        `gorm:"not null" json:"harga"`
	Stok          int        `gorm:"not null" json:"stok"`
	Kode          string     `gorm:"varchar(255);unique;not null" json:"kode"`
	Perusahaan_ID uuid.UUID  `gorm:"type:uuid;not null" json:"perusahaan_id"`
	Perusahaan    Perusahaan `gorm:"foreignKey:Perusahaan_ID;constraint:OnDelete:CASCADE"`
}
