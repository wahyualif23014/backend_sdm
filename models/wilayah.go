package models

// Struktur untuk output JSON ke Flutter
type WilayahResponse struct {
	Kode        string  `json:"id"`          // Kode Desa (misal: 35.01.01.2001)
	Kabupaten   string  `json:"kabupaten"`   // Nama Kabupaten
	Kecamatan   string  `json:"kecamatan"`   // Nama Kecamatan
	NamaDesa    string  `json:"namaDesa"`    // Nama Desa
	Latitude    float64 `json:"latitude"`    // Lat
	Longitude   float64 `json:"longitude"`   // Lng
	UpdatedBy   string  `json:"updatedBy"`   // ID Anggota / Nama
	LastUpdated string  `json:"lastUpdated"` // Waktu Transaksi
}
