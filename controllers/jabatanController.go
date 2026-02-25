package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wahyualif23014/backendGO/initializers"
	"github.com/wahyualif23014/backendGO/models"
)

// CreateJabatan: Menambah nama jabatan baru
func CreateJabatan(c *gin.Context) {
	var input struct {
		NamaJabatan string `json:"nama_jabatan" binding:"required"`
		IdAnggota   *int   `json:"id_anggota"` // Optional sesuai data di DB
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama jabatan wajib diisi"})
		return
	}

	jabatan := models.Jabatan{
		NamaJabatan:  input.NamaJabatan,
		DeleteStatus: "2", // Status Aktif
		IdAnggota:    input.IdAnggota,
	}

	// Operasi Create
	if err := initializers.DB.Create(&jabatan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data ke database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Jabatan berhasil ditambahkan",
		"data":    jabatan,
	})
}

// GetJabatan: Mengambil list jabatan aktif
func GetJabatan(c *gin.Context) {
	var jabatans []models.Jabatan
	// Filter deletestatus = '2' (Sesuai logic soft delete Anda)
	if err := initializers.DB.Where("deletestatus = ?", "2").Order("idjabatan DESC").Find(&jabatans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": jabatans})
}

// UpdateJabatan: Memperbarui nama jabatan berdasarkan idjabatan
func UpdateJabatan(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		NamaJabatan string `json:"nama_jabatan" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama jabatan wajib diisi"})
		return
	}

	var jabatan models.Jabatan
	// Mencari data yang aktif (status 2)
	if err := initializers.DB.Where("idjabatan = ? AND deletestatus = ?", id, "2").First(&jabatan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jabatan tidak ditemukan atau sudah dihapus"})
		return
	}

	// Update field namajabatan
	if err := initializers.DB.Model(&jabatan).Update("namajabatan", input.NamaJabatan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Jabatan berhasil diperbarui",
		"data":    jabatan,
	})
}

// DeleteJabatan: Soft delete (mengubah deletestatus menjadi '1')
func DeleteJabatan(c *gin.Context) {
	id := c.Param("id")

	var jabatan models.Jabatan
	// Cari data berdasarkan primary key idjabatan
	if err := initializers.DB.Where("idjabatan = ?", id).First(&jabatan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jabatan tidak ditemukan"})
		return
	}

	// Update deletestatus ke '1' (Soft Delete)
	if err := initializers.DB.Model(&jabatan).Update("deletestatus", "1").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Jabatan berhasil dihapus"})
}