package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wahyualif23014/backendGO/initializers"
	"github.com/wahyualif23014/backendGO/models"
)

type FilterOptions struct {
	Polres     []string `json:"polres"`
	Polsek     []string `json:"polsek"`
	JenisLahan []string `json:"jenis_lahan"`
	Komoditas  []string `json:"komoditas"`
}

// ========================================================
// 1. GET FILTER OPTIONS
// ========================================================
func GetKelolaFilterOptions(c *gin.Context) {

	var options FilterOptions
	selectedPolres := c.Query("polres")

	initializers.DB.Table("lahan").
		Joins("LEFT JOIN wilayah w_kab ON w_kab.kode = SUBSTR(lahan.idwilayah,1,5)").
		Where("w_kab.nama IS NOT NULL").
		Distinct("CONCAT('POLRES ', UPPER(w_kab.nama))").
		Pluck("CONCAT('POLRES ', UPPER(w_kab.nama))", &options.Polres)

	if selectedPolres != "" {
		namaKab := strings.TrimSpace(strings.TrimPrefix(selectedPolres, "POLRES "))

		initializers.DB.Table("lahan").
			Joins("LEFT JOIN wilayah w_kab ON w_kab.kode = SUBSTR(lahan.idwilayah,1,5)").
			Joins("LEFT JOIN wilayah w_kec ON w_kec.kode = SUBSTR(lahan.idwilayah,1,8)").
			Where("UPPER(w_kab.nama)=? AND w_kec.nama IS NOT NULL", namaKab).
			Distinct("CONCAT('POLSEK ', UPPER(w_kec.nama))").
			Pluck("CONCAT('POLSEK ', UPPER(w_kec.nama))", &options.Polsek)
	}

	options.JenisLahan = []string{"POKTAN BINAAN POLRI", "MASYARAKAT BINAAN"}
	options.Komoditas = []string{"JAGUNG", "PADI", "KEDELAI", "BAWANG MERAH", "CABAI"}

	c.JSON(http.StatusOK, options)
}

// ========================================================
// 2. GET SUMMARY (FIXED FINAL)
// ========================================================
func GetKelolaSummary(c *gin.Context) {

	var summary models.KelolaLahanSummary

	// Total Potensi
	initializers.DB.Table("lahan").
		Select("COALESCE(SUM(luaslahan),0)").
		Scan(&summary.TotalPotensiLahan)

	// Total Tanam
	initializers.DB.Table("tanam").
		Select("COALESCE(SUM(luastanam),0)").
		Scan(&summary.TotalTanamLahan)

	// Total Panen (Ha & Ton)
	initializers.DB.Table("panen").
		Select(`
			COALESCE(SUM(luaspanen),0),
			COALESCE(SUM(totalpanen),0)
		`).
		Row().
		Scan(
			&summary.TotalPanenLahanHa,
			&summary.TotalPanenLahanTon,
		)

	// Total Serapan dari distribusi
	initializers.DB.Table("distribusi").
		Select("COALESCE(SUM(totaldistribusi),0)").
		Scan(&summary.TotalSerapanTon)

	c.JSON(http.StatusOK, summary)
}

// ========================================================
// 3. GET LIST (FINAL STABLE)
// ========================================================
func GetKelolaList(c *gin.Context) {

	var result []models.KelolaLahanItem

	search := c.Query("search")
	polres := c.Query("polres")
	polsek := c.Query("polsek")
	jenisLahan := c.Query("jenis_lahan")

	query := initializers.DB.Table("lahan").
		Select(`
			lahan.idlahan as id,

			CONCAT('KAB. ', UPPER(w_kab.nama),
			       ' KEC. ', UPPER(w_kec.nama),
			       ' DESA ', UPPER(w_desa.nama)) as region_group,

			UPPER(lahan.alamat) as sub_region_group,

			-- Polisi Penggerak
			COALESCE(NULLIF(lahan.cp,''), '-') as police_name,
			COALESCE(NULLIF(lahan.hp,''), '-') as police_phone,

			-- Penanggung Jawab
			COALESCE(NULLIF(lahan.cppolisi,''), '-') as pic_name,
			COALESCE(NULLIF(lahan.hppolisi,''), '-') as pic_phone,

			COALESCE(lahan.luaslahan,0) as land_area,
			COALESCE(t.total_tanam,0) as luas_tanam,

			COALESCE(DATE_FORMAT(t.est_panen,'%d/%m/%Y'), '-') as est_panen,

			COALESCE(p.total_panen_ha,0) as luas_panen,
			COALESCE(p.total_panen_ton,0) as berat_panen,
			COALESCE(d.total_serapan,0) as serapan,

			CASE
				WHEN lahan.validoleh IS NOT NULL THEN true
				ELSE false
			END as is_validated,

			CASE
				WHEN lahan.validoleh IS NOT NULL THEN 'VALIDATED'
				ELSE 'PENDING'
			END as status
		`).
		Joins("LEFT JOIN wilayah w_desa ON w_desa.kode = lahan.idwilayah").
		Joins("LEFT JOIN wilayah w_kec ON w_kec.kode = SUBSTR(lahan.idwilayah,1,8)").
		Joins("LEFT JOIN wilayah w_kab ON w_kab.kode = SUBSTR(lahan.idwilayah,1,5)").

		// TANAM
		Joins(`
			LEFT JOIN (
				SELECT idlahan,
					   SUM(luastanam) as total_tanam,
					   MAX(estawalpanen) as est_panen
				FROM tanam
				GROUP BY idlahan
			) t ON t.idlahan = lahan.idlahan
		`).

		// PANEN
		Joins(`
			LEFT JOIN (
				SELECT idlahan,
					   SUM(luaspanen) as total_panen_ha,
					   SUM(totalpanen) as total_panen_ton
				FROM panen
				GROUP BY idlahan
			) p ON p.idlahan = lahan.idlahan
		`).

		// DISTRIBUSI
		Joins(`
			LEFT JOIN (
				SELECT idlahan,
					   SUM(totaldistribusi) as total_serapan
				FROM distribusi
				GROUP BY idlahan
			) d ON d.idlahan = lahan.idlahan
		`)

	// ================= FILTER =================

	if search != "" {
		s := "%" + strings.ToUpper(search) + "%"
		query = query.Where(
			"UPPER(lahan.alamat) LIKE ? OR UPPER(w_desa.nama) LIKE ?",
			s, s,
		)
	}

	if polres != "" {
		kab := strings.TrimSpace(strings.TrimPrefix(polres, "POLRES "))
		query = query.Where("UPPER(w_kab.nama) LIKE ?", "%"+kab+"%")
	}

	if polsek != "" {
		kec := strings.TrimSpace(strings.TrimPrefix(polsek, "POLSEK "))
		query = query.Where("UPPER(w_kec.nama) LIKE ?", "%"+kec+"%")
	}

	if jenisLahan != "" {
		if jenisLahan == "POKTAN BINAAN POLRI" {
			query = query.Where("lahan.idjenislahan = 2")
		} else {
			query = query.Where("lahan.idjenislahan != 2")
		}
	}

	if err := query.Order("lahan.idwilayah ASC").Scan(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range result {
		if result[i].IsValidated {
			result[i].StatusColor = "#4CAF50"
		} else {
			result[i].StatusColor = "#FF9800"
		}
	}

	c.JSON(http.StatusOK, result)
}
