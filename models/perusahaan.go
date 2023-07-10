package models

import uuid "github.com/satori/go.uuid"

type Perusahaan struct {
	ID      uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Nama    string    `gorm:"varchar(255);not null" json:"nama"`
	Alamat  string    `gorm:"varchar(255);not null" json:"alamat"`
	No_Telp string    `gorm:"varchar(255);not null" json:"no_telp"`
	Kode    string    `gorm:"varchar(3);not null" json:"kode"`
}
