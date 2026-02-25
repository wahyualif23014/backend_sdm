package models

type KelolaLahanSummary struct {
	TotalPotensiLahan  float64 `json:"total_potensi"`
	TotalTanamLahan    float64 `json:"total_tanam"`
	TotalPanenLahanHa  float64 `json:"total_panen_ha"`
	TotalPanenLahanTon float64 `json:"total_panen_ton"`
	TotalSerapanTon    float64 `json:"total_serapan"`
}

type KelolaLahanItem struct {
	ID             string `json:"id"`
	RegionGroup    string `json:"region_group"`
	SubRegionGroup string `json:"sub_region_group"`

	// 1. Penanggung Jawab
	PicName  string `json:"pic_name"`
	PicPhone string `json:"pic_phone"`

	// 2. Luas (Ha)
	LandArea float64 `json:"land_area"`

	// 3. Tanam (Ha)
	LuasTanam float64 `json:"luas_tanam"`

	// 4. Est. Panen
	EstPanen string `json:"est_panen"`

	// 5 & 6. Panen (Ha & Ton)
	LuasPanen  float64 `json:"luas_panen"`
	BeratPanen float64 `json:"berat_panen"`

	// 7. Serapan (Ton)
	Serapan float64 `json:"serapan"`

	// 8. Polisi Penggerak
	PoliceName  string `json:"police_name"`
	PolicePhone string `json:"police_phone"`

	// 9. Validasi
	IsValidated   bool   `json:"is_validated"`
	Status        string `json:"status"`
	StatusColor   string `json:"status_color"`
	KategoriLahan string `json:"kategori_lahan"`
}
