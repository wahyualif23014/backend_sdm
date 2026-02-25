package models

import (
	"time"
)

type Jabatan struct {
	ID              uint64    `gorm:"column:idjabatan;primaryKey;autoIncrement" json:"id"`
	NamaJabatan     string    `gorm:"column:namajabatan;size:100" json:"nama_jabatan"`
	DeleteStatus    string    `gorm:"column:deletestatus;type:enum('1','2');default:'2'" json:"-"`
	IdAnggota       *int      `gorm:"column:idanggota" json:"id_anggota"` // Gunakan pointer karena bisa NULL
	DateTransaction time.Time `gorm:"column:datetransaction;autoCreateTime" json:"created_at"`
}

func (Jabatan) TableName() string {
	return "jabatan"
}