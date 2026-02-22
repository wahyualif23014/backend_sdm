package models

import "time"

// Sesuaikan struct dengan tabel 'komoditi' di SQL
type Komoditi struct {
	IDKomoditi      uint64    `gorm:"column:idkomoditi;primaryKey;autoIncrement"`
	JenisKomoditi   string    `gorm:"column:jeniskomoditi"` // Kolom ini ada di SQL
	NamaKomoditi    string    `gorm:"column:namakomoditi"`
	IDAnggota       uint64    `gorm:"column:idanggota"`
	DeleteStatus    string    `gorm:"column:deletestatus;default:'2'"`
	DateTransaction time.Time `gorm:"column:datetransaction"`
}

// Pastikan TableName mengembalikan 'komoditi' (bukan komoditas)
func (Komoditi) TableName() string {
	return "komoditi"
}

// Struct Response (Sesuaikan dengan kebutuhan Flutter)
type CategoryResponse struct {
	ID    string   `json:"id"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}

type CommodityItemResponse struct {
	ID         string `json:"id"`
	CategoryID string `json:"categoryId"`
	Name       string `json:"name"`
	IsSelected bool   `json:"isSelected"`
}
