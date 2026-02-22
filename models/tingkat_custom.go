package models

type KesatuanDetail struct {
	// Map kolom SQL ke Field Struct pakai tag gorm:"column:..."
	Kode        string `json:"kode_kesatuan" gorm:"column:kode"`
	NamaSatuan  string `json:"nama_satuan" gorm:"column:nama_satuan"`
	NamaPejabat string `json:"nama_pejabat" gorm:"column:nama_pejabat"`
	Jabatan     string `json:"jabatan_pejabat" gorm:"column:jabatan"`
	NoHP        string `json:"no_hp" gorm:"column:no_hp"`

	// Field ini diisi oleh Logic Go (Manipulasi String), bukan langsung dari SQL
	// Maka WAJIB pakai gorm:"-"
	Wilayah string `json:"wilayah" gorm:"-"`

	// Field Internal
	KodeInduk string `json:"-" gorm:"-"`

	// Field Hierarki (Array)
	// WAJIB pakai gorm:"-" agar GORM tidak error saat Scan
	TotalPolsek  int              `json:"total_polsek" gorm:"-"`
	DaftarPolsek []KesatuanDetail `json:"daftar_polsek,omitempty" gorm:"-"`
}
