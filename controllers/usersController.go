package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wahyualif23014/backendGO/initializers"
	"github.com/wahyualif23014/backendGO/models"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserInput struct {
	NamaLengkap string `json:"nama_lengkap" binding:"required"`
	IDTugas     string `json:"id_tugas" binding:"required"`
	Username    string `json:"username" binding:"required"`
	JabatanID   uint64 `json:"id_jabatan" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Role        string `json:"role" binding:"required"`
	NoTelp      string `json:"no_telp"`
}

type UpdateUserInput struct {
	NamaLengkap string  `json:"nama_lengkap"`
	NoTelp      string  `json:"no_telp"`
	IDTugas     string  `json:"id_tugas"`
	IDJabatan   *uint64 `json:"id_jabatan"`
	Role        string  `json:"role"`
}

func CreateUser(c *gin.Context) {
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 10)

	userValue, _ := c.Get("user")
	adminUser := userValue.(models.User)

	user := models.User{
		NamaLengkap:     input.NamaLengkap,
		Username:        input.Username,
		KataSandi:       string(hash),
		IDTugas:         input.IDTugas,
		IDJabatan:       &input.JabatanID,
		Role:            input.Role,
		NoTelp:          input.NoTelp,
		IDPengguna:      adminUser.ID,
		DeleteStatus:    models.StatusActive,
		DateTransaction: time.Now(),
	}

	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User berhasil didaftarkan", "data": user})
}

func GetUsers(c *gin.Context) {
	var users []models.User

	// Implementasi Preload untuk mengambil detail Jabatan & Tingkat
	err := initializers.DB.
		Preload("Jabatan").
		Preload("TingkatDetail"). 
		Where("deletestatus = ?", models.StatusActive).
		Find(&users).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	err := initializers.DB.
		Preload("Jabatan").
		Preload("TingkatDetail").
		Where("idanggota = ? AND deletestatus = ?", id, models.StatusActive).
		First(&user).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := initializers.DB.Where("idanggota = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	updates := make(map[string]interface{})
	if input.NamaLengkap != "" { updates["nama"] = input.NamaLengkap }
	if input.NoTelp != ""      { updates["hp"] = input.NoTelp }
	if input.IDTugas != ""      { updates["idtugas"] = input.IDTugas }
	if input.IDJabatan != nil   { updates["idjabatan"] = *input.IDJabatan }
	if input.Role != ""         { updates["statusadmin"] = input.Role }

	initializers.DB.Model(&user).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "data": user})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	result := initializers.DB.Model(&models.User{}).
		Where("idanggota = ?", id).
		Update("deletestatus", models.StatusDeleted)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func GetProfile(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"data": user})
}