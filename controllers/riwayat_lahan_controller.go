package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wahyualif23014/backendGO/initializers"
	"github.com/wahyualif23014/backendGO/models"
)

func GetRiwayatSummary(c *gin.Context) {
	var summary models.RiwayatLahanSummary
	// Ambil data dari database (Contoh query sum)
	initializers.DB.Table("lahan").Select("COALESCE(SUM(luaslahan), 0)").Scan(&summary.TotalPotensiLahan)
	// Isi field lainnya sesuai kebutuhan business logic
	c.JSON(http.StatusOK, summary)
}

func GetRiwayatList(c *gin.Context) {
	var result []models.RiwayatLahanItem
	search := c.Query("search")

	query := initializers.DB.Table("lahan").
		Select(`
			lahan.idlahan as id,
			CONCAT('KAB. ', UPPER(w_kab.nama), ' KEC. ', UPPER(w_kec.nama), ' DESA ', UPPER(w_desa.nama)) as region_group,
			UPPER(lahan.alamat) as sub_region_group,
			lahan.cppolisi as police_name,
			lahan.hppolisi as police_phone,
			lahan.cp as pic_name,
			lahan.hp as pic_phone,
			lahan.luaslahan as land_area,
			'POKTAN BINAAN POLRI' as land_category,
			'SELESAI PANEN' as status,
			'#4CAF50' as status_color
		`).
		Joins("LEFT JOIN wilayah w_desa ON w_desa.kode = lahan.idwilayah").
		Joins("LEFT JOIN wilayah w_kec ON w_kec.kode = SUBSTR(lahan.idwilayah, 1, 8)").
		Joins("LEFT JOIN wilayah w_kab ON w_kab.kode = SUBSTR(lahan.idwilayah, 1, 5)")

	if search != "" {
		s := "%" + strings.ToUpper(search) + "%"
		query = query.Where("lahan.alamat LIKE ? OR w_desa.nama LIKE ?", s, s)
	}

	query.Scan(&result)
	c.JSON(http.StatusOK, result)
}

// GET FILTER OPTIONS RIWAYAT
func GetRiwayatFilterOptions(c *gin.Context) {
	var options FilterOptions // Gunakan struct yang sama dengan Kelola Lahan
	selectedPolres := c.Query("polres")

	// Ambil list Polres yang ada di data riwayat (tabel lahan)
	initializers.DB.Table("lahan").
		Joins("LEFT JOIN wilayah w_kab ON w_kab.kode = SUBSTR(lahan.idwilayah, 1, 5)").
		Where("w_kab.nama IS NOT NULL").
		Distinct("CONCAT('POLRES ', UPPER(w_kab.nama))").
		Pluck("CONCAT('POLRES ', UPPER(w_kab.nama))", &options.Polres)

	if selectedPolres != "" {
		namaKab := strings.TrimSpace(strings.TrimPrefix(selectedPolres, "POLRES "))
		initializers.DB.Table("lahan").
			Joins("LEFT JOIN wilayah w_kab ON w_kab.kode = SUBSTR(lahan.idwilayah, 1, 5)").
			Joins("LEFT JOIN wilayah w_kec ON w_kec.kode = SUBSTR(lahan.idwilayah, 1, 8)").
			Where("UPPER(w_kab.nama) = ?", namaKab).
			Distinct("CONCAT('POLSEK ', UPPER(w_kec.nama))").
			Pluck("CONCAT('POLSEK ', UPPER(w_kec.nama))", &options.Polsek)
	}

	options.JenisLahan = []string{"POKTAN BINAAN POLRI", "MASYARAKAT BINAAN"}
	options.Komoditas = []string{"JAGUNG", "PADI", "KEDELAI"}

	c.JSON(http.StatusOK, options)
}
