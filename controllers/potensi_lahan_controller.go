package controllers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wahyualif23014/backendGO/initializers"
	"github.com/wahyualif23014/backendGO/models"
)

func GetImageFromDB(c *gin.Context) {
	filenameRaw := c.Param("filename")
	filename, err := url.QueryUnescape(filenameRaw)
	if err != nil {
		filename = filenameRaw
	}

	var lahan models.PotensiLahan
	result := initializers.DB.Table("lahan").Where("dokumentasi = ?", filename).First(&lahan)

	if result.Error != nil {
		fmt.Printf("[ERROR] Data tidak ditemukan untuk file: %s\n", filename)
		c.JSON(http.StatusNotFound, gin.H{"error": "Data SQL tidak ditemukan"})
		return
	}

	base64String := lahan.FotoLahan
	if base64String == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kolom SQL kosong"})
		return
	}

	if strings.Contains(base64String, ",") {
		base64String = strings.Split(base64String, ",")[1]
	}

	base64String = strings.TrimSpace(base64String)
	base64String = strings.ReplaceAll(base64String, "\n", "")
	base64String = strings.ReplaceAll(base64String, "\r", "")

	imageBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		imageBytes, err = base64.RawStdEncoding.DecodeString(base64String)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Base64 Rusak", "debug": err.Error()})
			return
		}
	}

	c.Data(http.StatusOK, "image/jpeg", imageBytes)
}

func GetPotensiLahan(c *gin.Context) {
	var daftarLahan []models.PotensiLahan

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	search := c.Query("search")
	polres := c.Query("polres")
	polsek := c.Query("polsek")

	db := initializers.DB.Table("lahan").
		Select(`
			lahan.*, 
			MAX(w_desa.nama) AS nama_desa, 
			MAX(w_kec.nama) AS nama_kecamatan, 
			MAX(w_kab.nama) AS nama_kabupaten,
			MAX(p.nama) AS nama_pemroses,
			MAX(v.nama) AS nama_validator,
			MAX(k.jeniskomoditi) AS jenis_komoditas_nama,
			MAX(k.namakomoditi) AS nama_komoditi_asli
		`).
		Joins("LEFT JOIN wilayah w_desa ON w_desa.kode = lahan.idwilayah").
		Joins("LEFT JOIN wilayah w_kec ON w_kec.kode = SUBSTR(lahan.idwilayah, 1, 8)").
		Joins("LEFT JOIN wilayah w_kab ON w_kab.kode = SUBSTR(lahan.idwilayah, 1, 5)").
		Joins("LEFT JOIN anggota p ON p.idanggota = lahan.editoleh").
		Joins("LEFT JOIN anggota v ON v.idanggota = lahan.validoleh").
		Joins("LEFT JOIN komoditi k ON k.idkomoditi = lahan.idkomoditi").
		Where("lahan.statuslahan IS NOT NULL AND lahan.statuslahan IN ('1', '2', '3', '4')")

	if search != "" {
		s := "%" + strings.ToLower(search) + "%"
		db = db.Where("LOWER(lahan.alamat) LIKE ? OR LOWER(lahan.poktan) LIKE ?", s, s)
	}

	if polres != "" {
		db = db.Where("w_kab.nama = ?", polres)
	}
	if polsek != "" {
		db = db.Where("w_kec.nama = ?", polsek)
	}

	// Gunakan Group By untuk memastikan satu baris per ID Lahan (mencegah duplikasi akibat JOIN)
	if err := db.Group("lahan.idlahan").Order("lahan.datetransaction DESC").Limit(limit).Offset(offset).Find(&daftarLahan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   daftarLahan,
	})
}

func GetFilterOptions(c *gin.Context) {
	var listPolres []string
	var listPolsek []string

	initializers.DB.Table("lahan").
		Select("DISTINCT w_kab.nama").
		Joins("LEFT JOIN wilayah w_kab ON w_kab.kode = SUBSTR(lahan.idwilayah, 1, 5)").
		Where("w_kab.nama IS NOT NULL").
		Pluck("w_kab.nama", &listPolres)

	initializers.DB.Table("lahan").
		Select("DISTINCT w_kec.nama").
		Joins("LEFT JOIN wilayah w_kec ON w_kec.kode = SUBSTR(lahan.idwilayah, 1, 8)").
		Where("w_kec.nama IS NOT NULL").
		Pluck("w_kec.nama", &listPolsek)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"polres": listPolres,
			"polsek": listPolsek,
		},
	})
}

func CreatePotensiLahan(c *gin.Context) {
	var input models.PotensiLahan
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	input.DateTransaction = time.Now()
	if err := initializers.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": input})
}

func UpdatePotensiLahan(c *gin.Context) {
	id := c.Param("id")
	var lahan models.PotensiLahan
	if err := initializers.DB.First(&lahan, "idlahan = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data tidak ditemukan"})
		return
	}
	var input models.PotensiLahan
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	initializers.DB.Model(&lahan).Updates(input)
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": lahan})
}

func DeletePotensiLahan(c *gin.Context) {
	id := c.Param("id")
	if err := initializers.DB.Delete(&models.PotensiLahan{}, "idlahan = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Data berhasil dihapus"})
}

func GetSummaryLahan(c *gin.Context) {
	type SummaryCategory struct {
		Title string  `json:"title"`
		Area  float64 `json:"area"`
		Count int64   `json:"count"`
	}

	var totals struct {
		TotalArea float64 `gorm:"column:total_area"`
		TotalLoc  int64   `gorm:"column:total_loc"`
	}

	dbFilter := "idwilayah IS NOT NULL AND statuslahan IN ('1', '2', '3', '4')"

	// Menggunakan subquery untuk memastikan SUM hanya menghitung baris unik
	initializers.DB.Table("lahan").
		Where(dbFilter).
		Select("COALESCE(SUM(luaslahan), 0) as total_area, COUNT(DISTINCT idlahan) as total_loc").
		Scan(&totals)

	var categories []SummaryCategory

	rows, err := initializers.DB.Table("lahan").
		Where(dbFilter).
		Select(`
			idjenislahan,
			COALESCE(SUM(luaslahan), 0) as area,
			COUNT(DISTINCT idlahan) as count
		`).
		Group("idjenislahan").
		Rows()

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id int
			var area float64
			var count int64
			rows.Scan(&id, &area, &count)

			title := "LAHAN LAINNYA"
			switch id {
			case 1:
				title = "PERHUTANAN SOSIAL"
			case 2:
				title = "POKTAN BINAAN POLRI"
			case 3:
				title = "MASYARAKAT BINAAN POLRI"
			case 4:
				title = "TUMPANG SARI"
			case 5:
				title = "MILIK POLRI"
			case 6:
				title = "LBS"
			case 7:
				title = "PESANTREN"
			}

			categories = append(categories, SummaryCategory{
				Title: title,
				Area:  area,
				Count: count,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"total_area":      totals.TotalArea,
			"total_locations": totals.TotalLoc,
			"categories":      categories,
		},
	})
}

func GetNoPotentialLahan(c *gin.Context) {
	var master struct {
		Kab  int64
		Kec  int64
		Desa int64
	}

	initializers.DB.Table("wilayah").Select(`
		SUM(CASE WHEN CHAR_LENGTH(kode) = 5 THEN 1 ELSE 0 END) as kab,
		SUM(CASE WHEN CHAR_LENGTH(kode) = 8 THEN 1 ELSE 0 END) as kec,
		SUM(CASE WHEN CHAR_LENGTH(kode) > 8 THEN 1 ELSE 0 END) as desa
	`).Scan(&master)

	var isi struct {
		Kab  int64
		Kec  int64
		Desa int64
	}

	dbFilter := "idwilayah IS NOT NULL AND statuslahan IN ('1', '2', '3', '4')"

	initializers.DB.Table("lahan").
		Where(dbFilter).
		Select(`
			COUNT(DISTINCT SUBSTR(idwilayah, 1, 5)) as kab,
			COUNT(DISTINCT SUBSTR(idwilayah, 1, 8)) as kec,
			COUNT(DISTINCT idwilayah) as desa
		`).Scan(&isi)

	polresKosong := master.Kab - isi.Kab
	polsekKosong := master.Kec - isi.Kec
	desaKosong := master.Desa - isi.Desa

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"total_empty_polres": polresKosong,
			"details": gin.H{
				"polsek":    polsekKosong,
				"kab_kota":  polresKosong,
				"kecamatan": polsekKosong,
				"kel_desa":  desaKosong,
			},
		},
	})
}
