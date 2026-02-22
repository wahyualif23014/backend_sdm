package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wahyualif23014/backendGO/initializers"
	"github.com/wahyualif23014/backendGO/models"
)

func GetTingkat(c *gin.Context) {
	var rawData []models.KesatuanDetail

	query := `
		SELECT 
			t.kode AS kode,
			t.nama AS nama_satuan,
			COALESCE(j.namajabatan, '-') AS jabatan,
			COALESCE(a.nama, 'Belum Ada Pejabat') AS nama_pejabat,
			COALESCE(a.hp, '-') AS no_hp
		FROM tingkat t
		LEFT JOIN (
			SELECT *, 
			ROW_NUMBER() OVER (
				PARTITION BY idtugas 
				ORDER BY FIELD(idjabatan, 1, 2, 3, 7, 8) ASC
			) as rank_jabatan
			FROM anggota 
			WHERE deletestatus != '1' 
		) a ON a.idtugas = t.kode AND a.rank_jabatan = 1
		LEFT JOIN jabatan j ON a.idjabatan = j.idjabatan
		ORDER BY t.kode ASC
	`

	if err := initializers.DB.Raw(query).Scan(&rawData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error: " + err.Error()})
		return
	}

	var polresList []models.KesatuanDetail
	anakMap := make(map[string][]models.KesatuanDetail)

	for i := range rawData {
		kode := rawData[i].Kode
		nama := rawData[i].NamaSatuan

		cleanName := nama
		cleanName = strings.Replace(cleanName, "POLSEK ", "", 1)
		cleanName = strings.Replace(cleanName, "POLRES ", "", 1)
		cleanName = strings.Replace(cleanName, "POLRESTA ", "", 1)
		cleanName = strings.Replace(cleanName, "POLRESTABES ", "", 1)
		cleanName = strings.Replace(cleanName, "POLDA ", "", 1)
		rawData[i].Wilayah = cleanName

		if len(kode) > 5 {
			indukID := kode[:5]
			rawData[i].KodeInduk = indukID
			anakMap[indukID] = append(anakMap[indukID], rawData[i])
		} else {
			polresList = append(polresList, rawData[i])
		}
	}

	for i := range polresList {
		idInduk := polresList[i].Kode
		if anak, ok := anakMap[idInduk]; ok {
			polresList[i].DaftarPolsek = anak
			polresList[i].TotalPolsek = len(anak)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   polresList,
	})
}