package models

import (
	"time"
)

const (
	RoleAdmin    string = "1"
	RoleOperator string = "2"
	RoleView     string = "3"
)

const (
	StatusDeleted string = "1"
	StatusActive  string = "2"
)

type User struct {
	ID          uint64 `gorm:"column:idanggota;primaryKey;autoIncrement" json:"id"`
	NamaLengkap string `gorm:"column:nama;size:100" json:"nama_lengkap"`
	NoTelp      string `gorm:"column:hp;size:20" json:"no_telp"`
	IDTugas     string `gorm:"column:idtugas;size:13;not null" json:"id_tugas"`

	Username        string    `gorm:"column:username;type:longtext" json:"nrp"`
	KataSandi       string    `gorm:"column:password;type:longtext" json:"-"`
	Role            string    `gorm:"column:statusadmin;type:enum('1','2','3');default:'3'" json:"role"`
	DeleteStatus    string    `gorm:"column:deletestatus;type:enum('1','2');default:'2'" json:"-"`
	DateTransaction time.Time `gorm:"column:datetransaction;autoCreateTime" json:"created_at"`
	IDPengguna      uint64    `gorm:"column:idpengguna;not null" json:"id_pengguna"`
	TingkatDetail   *Tingkat  `gorm:"foreignKey:IDTugas;references:Kode" json:"tingkat_detail,omitempty"`

	// Relasi Belongs-To ke Jabatan
	IDJabatan *uint64  `gorm:"column:idjabatan" json:"id_jabatan"`
	Jabatan   *Jabatan `gorm:"foreignKey:IDJabatan;references:ID" json:"jabatan,omitempty"`
}

func (User) TableName() string {
	return "anggota"
}
