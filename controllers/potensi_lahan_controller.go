// package controllers

// import (
// 	"net/http"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/wahyualif23014/backendGO/initializers"
// 	"github.com/wahyualif23014/backendGO/models"
// )

// // GET: Ambil Data Utama dengan Optimasi Query
// func GetPotensiLahan(c *gin.Context) {
// 	var daftarLahan []models.PotensiLahan

// 	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
// 	offset := (page - 1) * limit

// 	search := c.Query("search")
// 	status := c.Query("status")
// 	polres := c.Query("polres")
// 	polsek := c.Query("polsek")
// 	jenisLahan := c.Query("jenis_lahan")

// 	// Optimasi: Gunakan SUBSTR untuk efisiensi index (jika database mendukung)
// 	// Gunakan Select eksplisit untuk menghindari ambiguitas kolom
// 	db := initializers.DB.Table("lahan").
// 		Select(`
// 			DISTINCT lahan.*,
// 			w_desa.nama AS nama_desa,
// 			w_kec.nama AS nama_kecamatan,
// 			w_kab.nama AS nama_kabupaten
// 		`).
// 		Joins("LEFT JOIN wilayah w_desa ON w_desa.kode = lahan.idwilayah").
// 		Joins("LEFT JOIN wilayah w_kec ON w_kec.kode = SUBSTR(lahan.idwilayah, 1, 8)").
// 		Joins("LEFT JOIN wilayah w_kab ON w_kab.kode = SUBSTR(lahan.idwilayah, 1, 5)")

// 	if search != "" {
// 		s := "%" + strings.ToLower(search) + "%"
// 		db = db.Where("LOWER(lahan.alamat) LIKE ? OR LOWER(lahan.poktan) LIKE ?", s, s)
// 	}
// 	if status != "" {
// 		db = db.Where("lahan.statuslahan = ?", status)
// 	}
// 	if polres != "" {
// 		db = db.Where("w_kab.nama = ?", polres)
// 	}
// 	if polsek != "" {
// 		db = db.Where("w_kec.nama = ?", polsek)
// 	}
// 	if jenisLahan != "" {
// 		val := 1
// 		if jenisLahan == "LADANG" {
// 			val = 2
// 		}
// 		db = db.Where("lahan.idjenislahan = ?", val)
// 	}

// 	// Gunakan Find untuk performa GORM yang lebih stabil dibanding Scan pada slice model
// 	if err := db.Order("lahan.datetransaction DESC").Limit(limit).Offset(offset).Find(&daftarLahan).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"status": "success",
// 		"data":   daftarLahan,
// 		"page":   page,
// 		"limit":  limit,
// 	})
// }

// // GET: Opsi Dropdown Filter
// func GetFilterOptions(c *gin.Context) {
// 	var listPolres []string
// 	var listPolsek []string

// 	initializers.DB.Table("lahan").
// 		Select("DISTINCT w_kab.nama").
// 		Joins("LEFT JOIN wilayah w_kab ON w_kab.kode = SUBSTR(lahan.idwilayah, 1, 5)").
// 		Where("w_kab.nama IS NOT NULL").
// 		Pluck("w_kab.nama", &listPolres)

// 	initializers.DB.Table("lahan").
// 		Select("DISTINCT w_kec.nama").
// 		Joins("LEFT JOIN wilayah w_kec ON w_kec.kode = SUBSTR(lahan.idwilayah, 1, 8)").
// 		Where("w_kec.nama IS NOT NULL").
// 		Pluck("w_kec.nama", &listPolsek)

// 	c.JSON(http.StatusOK, gin.H{
// 		"status": "success",
// 		"data": gin.H{
// 			"polres": listPolres,
// 			"polsek": listPolsek,
// 		},
// 	})
// }

// // POST: Create
// func CreatePotensiLahan(c *gin.Context) {
// 	var input models.PotensiLahan
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}
// 	input.DateTransaction = time.Now()
// 	if err := initializers.DB.Create(&input).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": input})
// }

// // PUT: Update
// func UpdatePotensiLahan(c *gin.Context) {
// 	id := c.Param("id")
// 	var lahan models.PotensiLahan
// 	if err := initializers.DB.First(&lahan, "idlahan = ?", id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data tidak ditemukan"})
// 		return
// 	}
// 	var input models.PotensiLahan
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}
// 	initializers.DB.Model(&lahan).Updates(input)
// 	c.JSON(http.StatusOK, gin.H{"status": "success", "data": lahan})
// }

// // DELETE: Hapus
// func DeletePotensiLahan(c *gin.Context) {
// 	id := c.Param("id")
// 	if err := initializers.DB.Delete(&models.PotensiLahan{}, "idlahan = ?", id).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Data berhasil dihapus"})
// }

// // GET: Summary Data
// func GetSummaryLahan(c *gin.Context) {
// 	type SummaryCategory struct {
// 		Title string  `json:"title"`
// 		Area  float64 `json:"area"`
// 		Count int64   `json:"count"`
// 	}

// 	var totalArea float64
// 	var totalLoc int64
// 	var categories []SummaryCategory

// 	initializers.DB.Table("lahan").Select("COALESCE(SUM(luaslahan), 0)").Row().Scan(&totalArea)
// 	initializers.DB.Table("lahan").Count(&totalLoc)

// 	var kabCount, kecCount, desaCount int64
// 	initializers.DB.Table("lahan").Select("COUNT(DISTINCT SUBSTR(idwilayah, 1, 5))").Scan(&kabCount)
// 	initializers.DB.Table("lahan").Select("COUNT(DISTINCT SUBSTR(idwilayah, 1, 8))").Scan(&kecCount)
// 	initializers.DB.Table("lahan").Select("COUNT(DISTINCT idwilayah)").Scan(&desaCount)

// 	rows, err := initializers.DB.Table("lahan").
// 		Select("idjenislahan, COALESCE(SUM(luaslahan), 0) as area, COUNT(idlahan) as count").
// 		Group("idjenislahan").Rows()

// 	if err == nil {
// 		defer rows.Close()
// 		for rows.Next() {
// 			var id int
// 			var area float64
// 			var count int64
// 			rows.Scan(&id, &area, &count)

// 			title := "LAHAN LAINNYA"
// 			switch id {
// 			case 1:
// 				title = "MILIK POLRI"
// 			case 2:
// 				title = "POKTAN BINAAN POLRI"
// 			case 3:
// 				title = "MASYARAKAT BINAAN POLRI"
// 			case 4:
// 				title = "TUMPANG SARI"
// 			case 5:
// 				title = "PERHUTANAN SOSIAL"
// 			case 6:
// 				title = "LBS"
// 			case 7:
// 				title = "PESANTREN"
// 			}
// 			categories = append(categories, SummaryCategory{Title: title, Area: area, Count: count})
// 		}
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"status": "success",
// 		"data": gin.H{
// 			"total_area":      totalArea,
// 			"total_locations": totalLoc,
// 			"categories":      categories,
// 			"admin_counts": gin.H{
// 				"polres":    kabCount,
// 				"polsek":    kecCount,
// 				"kab_kota":  kabCount,
// 				"kecamatan": kecCount,
// 				"kel_desa":  desaCount,
// 			},
// 		},
// 	})
// }

// // GET: Wilayah Kosong
// func GetNoPotentialLahan(c *gin.Context) {
// 	var totalKabMaster, totalKecMaster, totalDesaMaster int64
// 	var isiKab, isiKec, isiDesa int64

// 	initializers.DB.Table("wilayah").Where("CHAR_LENGTH(kode) = 5").Count(&totalKabMaster)
// 	initializers.DB.Table("wilayah").Where("CHAR_LENGTH(kode) = 8").Count(&totalKecMaster)
// 	initializers.DB.Table("wilayah").Where("CHAR_LENGTH(kode) > 8").Count(&totalDesaMaster)

// 	initializers.DB.Table("lahan").Select("COUNT(DISTINCT SUBSTR(idwilayah, 1, 5))").Scan(&isiKab)
// 	initializers.DB.Table("lahan").Select("COUNT(DISTINCT SUBSTR(idwilayah, 1, 8))").Scan(&isiKec)
// 	initializers.DB.Table("lahan").Select("COUNT(DISTINCT idwilayah)").Scan(&isiDesa)

//		c.JSON(http.StatusOK, gin.H{
//			"status": "success",
//			"data": gin.H{
//				"total_empty_polres": totalKabMaster - isiKab,
//				"details": gin.H{
//					"polsek":    totalKecMaster - isiKec,
//					"kab_kota":  totalKabMaster - isiKab,
//					"kecamatan": totalKecMaster - isiKec,
//					"kel_desa":  totalDesaMaster - isiDesa,
//				},
//			},
//		})
//	}
package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wahyualif23014/backendGO/initializers"
	"github.com/wahyualif23014/backendGO/models"
)

func GetPotensiLahan(c *gin.Context) {
	var lahan []models.PotensiLahan
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	query := initializers.DB.Table("lahan").
		Select("lahan.*, w_desa.nama AS nama_desa, w_kec.nama AS nama_kecamatan, w_kab.nama AS nama_kabupaten").
		Joins("LEFT JOIN wilayah w_desa ON w_desa.kode = lahan.idwilayah").
		Joins("LEFT JOIN wilayah w_kec ON w_kec.kode = LEFT(lahan.idwilayah, 8)").
		Joins("LEFT JOIN wilayah w_kab ON w_kab.kode = LEFT(lahan.idwilayah, 5)").
		Where("lahan.deletestatus = ?", "2")

	query.Limit(limit).Offset(offset).Find(&lahan)

	c.JSON(http.StatusOK, gin.H{"data": lahan})
}

func CreatePotensiLahan(c *gin.Context) {
	var input models.PotensiLahan
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userValue, _ := c.Get("user")
	user := userValue.(models.User)

	// Inisialisasi metadata
	input.DateTransaction = time.Now()
	input.IDAnggota = user.ID
	input.IDTingkat = user.IDTugas
	input.DeleteStatus = "2"

	// PROTEKSI ENUM: Jika status_lahan tidak dikirim atau kosong, set default ke '1'
	if input.StatusLahan == "" {
		input.StatusLahan = "1"
	}

	if err := initializers.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Data lahan berhasil disimpan",
		"data":    input,
	})
}