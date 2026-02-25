package models

import (
	"time"
)

type PotensiLahan struct {
	// ID Utama
	ID        uint64 `gorm:"primaryKey;column:idlahan" json:"id"`
	IDTingkat string `gorm:"column:idtingkat" json:"id_tingkat"`
	IDWilayah string `gorm:"column:idwilayah" json:"id_wilayah"`

	// --- FIELD HASIL JOIN (Virtual Fields) ---
	// Field ini diisi menggunakan JOIN di controller
	NamaKabupaten     string `gorm:"column:nama_kabupaten;->" json:"nama_kabupaten"`
	NamaKecamatan     string `gorm:"column:nama_kecamatan;->" json:"nama_kecamatan"`
	NamaDesa          string `gorm:"column:nama_desa;->" json:"nama_desa"`
	NamaPemroses      string `gorm:"column:nama_pemroses;->" json:"nama_pemroses"`
	NamaValidator     string `gorm:"column:nama_validator;->" json:"nama_validator"`
	JenisKomoditiNama string `gorm:"column:jenis_komoditas_nama;->" json:"jenis_komoditas_nama"`
	NamaKomoditiAsli  string `gorm:"column:nama_komoditi_asli;->" json:"nama_komoditi_asli"`

	// --- DATA LAHAN ---
	IDJenisLahan int     `gorm:"column:idjenislahan" json:"id_jenis_lahan"`
	Alamat       string  `gorm:"column:alamat" json:"alamat_lahan"`
	LuasLahan    float64 `gorm:"column:luaslahan" json:"luas_lahan"`
	NamaPoktan   string  `gorm:"column:poktan" json:"poktan"`

	// --- CONTACT PERSON ---
	CPName      string `gorm:"column:cp" json:"pic_name"`
	CPPhone     string `gorm:"column:hp" json:"pic_phone"`
	PolisiName  string `gorm:"column:cppolisi" json:"police_name"`
	PolisiPhone string `gorm:"column:hppolisi" json:"police_phone"`

	// --- STATISTIK & KETERANGAN ---
	JumlahPoktan   int    `gorm:"column:jumlah_poktan;default:0" json:"jumlah_poktan"`
	JumlahPetani   int    `gorm:"column:jumlah_petani;default:0" json:"jumlah_petani"`
	KeteranganLain string `gorm:"column:keterangan_lain" json:"keterangan_lain"`

	// --- MEDIA & VALIDASI ---
	// FotoLahan digunakan untuk menyimpan string base64 dari DB
	FotoLahan   string `gorm:"column:foto_lahan" json:"foto_base64,omitempty"`
	Foto        string `gorm:"column:dokumentasi" json:"foto_lahan"`
	StatusLahan string `gorm:"column:statuslahan" json:"status_validasi"`
	IDKomoditi  int    `gorm:"column:idkomoditi" json:"id_komoditi"`

	// --- AUDIT TRAIL ---
	DateTransaction time.Time `gorm:"column:datetransaction" json:"tgl_proses"`
	EditOleh        string    `gorm:"column:editoleh" json:"edit_oleh"`
	ValidOleh       string    `gorm:"column:validoleh" json:"valid_oleh"`
	TglValidasi     string    `gorm:"column:tgl_validasi" json:"tgl_validasi"`

	// Field pendukung JSON (Opsional jika ingin mapping nama lama)
	DiprosesOleh   string `gorm:"column:diproses_oleh" json:"diproses_oleh"`
	DivalidasiOleh string `gorm:"column:divalidasi_oleh" json:"divalidasi_oleh"`
}

func (PotensiLahan) TableName() string {
	return "lahan"
}
