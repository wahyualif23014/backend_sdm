package models

// Untuk Kotak Statistik Atas
type RiwayatLahanSummary struct {
	TotalPotensiLahan  float64 `json:"total_potensi"`
	TotalTanamLahan    float64 `json:"total_tanam"`
	TotalPanenLahanHa  float64 `json:"total_panen_ha"`
	TotalPanenLahanTon float64 `json:"total_panen_ton"`
	TotalSerapanTon    float64 `json:"total_serapan"`
}

// Untuk Baris Data List
type RiwayatLahanItem struct {
	ID             string  `json:"id"`
	RegionGroup    string  `json:"region_group"`
	SubRegionGroup string  `json:"sub_region_group"`
	PoliceName     string  `json:"police_name"`
	PolicePhone    string  `json:"police_phone"`
	PicName        string  `json:"pic_name"`
	PicPhone       string  `json:"pic_phone"`
	LandArea       float64 `json:"land_area"`
	LandCategory   string  `json:"land_category"`
	Status         string  `json:"status"`
	StatusColor    string  `json:"status_color"`
	CreatedAt      string  `json:"created_at"` // Tanggal riwayat
}