package initializers

import (
	"log"
	// "github.com/wahyualif23014/backendGO/models" // Bisa dikomentari jika tidak dipakai
)

func SyncDatabase() {
	if DB == nil {
		log.Fatal("Database connection failed. Cannot sync.")
	}

	// ------------------------------------------------------------------
	// PENTING: MATIKAN AUTO MIGRATE UNTUK DATABASE 'presisi'
	// ------------------------------------------------------------------
	// Karena database 'presisi' sudah memiliki struktur tabel yang paten,
	// kita tidak boleh membiarkan GORM mengutak-atik strukturnya.
	// Jika kode ini dijalankan, GORM akan mencoba menambahkan Primary Key
	// ke tabel yang sudah punya Primary Key -> ERROR 1068.
	
	/*
	err := DB.AutoMigrate(
		&models.User{},    
		&models.Wilayah{}, 
		&models.Polres{},  
		&models.Polsek{},  
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	*/

	log.Println("Database migration SKIPPED (Using existing 'presisi' schema).")
}

