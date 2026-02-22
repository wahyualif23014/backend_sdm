package controllers

import (
	"fmt" // <--- JANGAN LUPA IMPORT INI
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wahyualif23014/backendGO/initializers"
	"github.com/wahyualif23014/backendGO/models"
)

// 1. GET KATEGORI & TOTAL STATISTIK
func GetCategories(c *gin.Context) {
	var categories []string
	var totalItems int64

	// DEBUG: Cek apakah model tersambung ke tabel yang benar
	fmt.Println("--- DEBUG GET CATEGORIES ---")

	// Hitung Total
	if err := initializers.DB.Model(&models.Komoditi{}).
		Where("deletestatus = ?", "2").
		Where("namakomoditi IS NOT NULL AND namakomoditi != ''").
		Count(&totalItems).Error; err != nil {

		fmt.Println("Error Count:", err) // Print Error jika ada
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung data tanaman"})
		return
	}

	// Print Hasil Hitungan ke Terminal
	fmt.Printf("Total Tanaman Ditemukan: %d\n", totalItems)

	// ... (Sisa kode query categories tetap sama) ...
	result := initializers.DB.Model(&models.Komoditi{}).
		Where("deletestatus = ?", "2").
		Where("jeniskomoditi IS NOT NULL AND jeniskomoditi != ''").
		Distinct().
		Pluck("jeniskomoditi", &categories)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data kategori"})
		return
	}

	var response []gin.H
	for i, cat := range categories {
		response = append(response, gin.H{
			"id":         strconv.Itoa(i + 1),
			"title":      cat,
			"imageAsset": "",
		})
	}

	if response == nil {
		response = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":        response,
		"total_items": totalItems,
	})
}

// ==========================================
// 2. GET COMMODITIES (Detail by Kind)
// ==========================================
func GetCommodities(c *gin.Context) {
	kind := c.Query("kind") // Ambil parameter ?kind=...
	var dbData []models.Komoditi

	// Query Dasar (Hanya data aktif)
	query := initializers.DB.Where("deletestatus = ?", "2")

	// Filter berdasarkan jenis jika ada parameter kind
	if kind != "" {
		query = query.Where("jeniskomoditi = ?", kind)
	}

	// Eksekusi Query
	if err := query.Find(&dbData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data detail"})
		return
	}

	var responseData []models.CommodityItemResponse
	for _, item := range dbData {
		responseData = append(responseData, models.CommodityItemResponse{
			ID:         strconv.FormatUint(item.IDKomoditi, 10),
			CategoryID: item.JenisKomoditi,
			Name:       item.NamaKomoditi,
			IsSelected: false,
		})
	}

	if responseData == nil {
		responseData = []models.CommodityItemResponse{}
	}

	c.JSON(http.StatusOK, gin.H{"data": responseData})
}

// ==========================================
// 3. CREATE COMMODITY (Tambah Data)
// ==========================================
func CreateCommodity(c *gin.Context) {
	var input struct {
		Name       string `json:"name" binding:"required"`
		CategoryID string `json:"categoryId"` // Ini menyimpan String Nama Jenis (ex: "Sayuran")
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newKomoditi := models.Komoditi{
		NamaKomoditi:    input.Name,
		JenisKomoditi:   input.CategoryID,
		DeleteStatus:    "2",
		IDAnggota:       1, // Default, bisa diganti user ID login
		DateTransaction: time.Now(),
	}

	if err := initializers.DB.Create(&newKomoditi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil", "data": newKomoditi})
}

// ==========================================
// 4. DELETE CATEGORY (Bulk Delete)
// ==========================================
func DeleteCategory(c *gin.Context) {
	var input struct {
		KindName string `json:"kindName" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama kategori diperlukan"})
		return
	}

	// Soft Delete (ubah deletestatus jadi '1') untuk SEMUA data dengan jenis tersebut
	result := initializers.DB.Model(&models.Komoditi{}).
		Where("jeniskomoditi = ?", input.KindName).
		Update("deletestatus", "1")

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil dihapus"})
}

// ==========================================
// 5. UPDATE COMMODITY (Edit Nama Tanaman)
// ==========================================
func UpdateCommodity(c *gin.Context) {
	var input struct {
		ID   string `json:"id" binding:"required"`
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update nama komoditi berdasarkan ID
	result := initializers.DB.Model(&models.Komoditi{}).
		Where("idkomoditi = ?", input.ID).
		Update("namakomoditi", input.Name)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diupdate"})
}

// ==========================================
// 6. DELETE COMMODITY ITEM (Hapus 1 Tanaman)
// ==========================================
func DeleteCommodityItem(c *gin.Context) {
	var input struct {
		ID string `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID diperlukan"})
		return
	}

	// Soft Delete satu item
	result := initializers.DB.Model(&models.Komoditi{}).
		Where("idkomoditi = ?", input.ID).
		Update("deletestatus", "1")

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item berhasil dihapus"})
}