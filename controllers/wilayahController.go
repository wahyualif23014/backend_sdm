// controllers/wilayahController.go
package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wahyualif23014/backendGO/initializers"
	"github.com/wahyualif23014/backendGO/models"
)

func GetWilayah(c *gin.Context) {
	results := []models.WilayahResponse{}
	query := `
		SELECT 
			d.kode,
			COALESCE(k.nama, '') AS kabupaten,
			COALESCE(c.nama, '') AS kecamatan,
			d.nama AS nama_desa,
			COALESCE(d.lat, 0) AS latitude,
			COALESCE(d.lng, 0) AS longitude,
			COALESCE(d.idanggota, '') AS updated_by,
			COALESCE(DATE_FORMAT(d.datetransaction, '%Y-%m-%d %H:%i:%s'), '') AS last_updated
		FROM wilayah d
		LEFT JOIN wilayah c ON c.kode = LEFT(d.kode, 8) 
		LEFT JOIN wilayah k ON k.kode = LEFT(d.kode, 5) 
		WHERE CHAR_LENGTH(d.kode) > 8 
		ORDER BY d.kode ASC
	`
	if err := initializers.DB.Raw(query).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func UpdateWilayah(c *gin.Context) {
	kode := c.Param("id")
	var body struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	// Ambil user dari middleware JWT
	userValue, _ := c.Get("user")
	currentUser := userValue.(models.User)

	result := initializers.DB.Table("wilayah").
		Where("kode = ?", kode).
		Updates(map[string]interface{}{
			"lat":             body.Latitude,
			"lng":             body.Longitude,
			"idanggota":       currentUser.ID,
			"datetransaction": time.Now(),
		})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil diperbarui"})
}
