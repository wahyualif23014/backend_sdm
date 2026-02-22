package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wahyualif23014/backendGO/initializers"
	"github.com/wahyualif23014/backendGO/models"
)

func GetJabatan(c *gin.Context) {
	var jabatans []models.Jabatan
	initializers.DB.Where("deletestatus = ?", "2").Find(&jabatans)
	c.JSON(http.StatusOK, gin.H{"data": jabatans})
}

func CreateJabatan(c *gin.Context) {
	var input struct {
		NamaJabatan string `json:"nama_jabatan" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jabatan := models.Jabatan{
		NamaJabatan:  input.NamaJabatan,
		DeleteStatus: "2",
	}

	initializers.DB.Create(&jabatan)
	c.JSON(http.StatusCreated, gin.H{"message": "Jabatan created", "data": jabatan})
}

func UpdateJabatan(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		NamaJabatan string `json:"nama_jabatan" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var jabatan models.Jabatan
	if err := initializers.DB.Where("idjabatan = ?", id).First(&jabatan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jabatan not found"})
		return
	}

	initializers.DB.Model(&jabatan).Update("namajabatan", input.NamaJabatan)
	c.JSON(http.StatusOK, gin.H{"message": "Jabatan updated", "data": jabatan})
}

func DeleteJabatan(c *gin.Context) {
	id := c.Param("id")
	result := initializers.DB.Model(&models.Jabatan{}).Where("idjabatan = ?", id).Update("deletestatus", "1")
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jabatan not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Jabatan deleted successfully"})
}