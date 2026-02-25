package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/wahyualif23014/backendGO/controllers"
	"github.com/wahyualif23014/backendGO/initializers"
	"github.com/wahyualif23014/backendGO/middleware"
	"github.com/wahyualif23014/backendGO/models"

	_ "github.com/wahyualif23014/backendGO/docs"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	r.GET("/uploads/:filename", controllers.GetImageFromDB)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "SIKAP PRESISI Backend v2.0 Online"})
	})

	r.POST("/login", controllers.Login)

	api := r.Group("/api")
	api.Use(middleware.RequireAuth)
	{
		// A. ADMIN ONLY
		admin := api.Group("/admin")
		admin.Use(middleware.RequireRoles(models.RoleAdmin))
		{
			admin.POST("/users", controllers.CreateUser)
			admin.GET("/users", controllers.GetUsers)
			admin.GET("/users/:id", controllers.GetUserByID)
			admin.PUT("/users/:id", controllers.UpdateUser)
			admin.DELETE("/users/:id", controllers.DeleteUser)
			admin.POST("/jabatan", controllers.CreateJabatan)
			admin.GET("/jabatan", controllers.GetJabatan)
			admin.PUT("/jabatan/:id", controllers.UpdateJabatan)
			admin.DELETE("/jabatan/:id", controllers.DeleteJabatan)
			admin.GET("/tingkat", controllers.GetTingkat)
			admin.PUT("/wilayah/:id", controllers.UpdateWilayah)
			admin.GET("/wilayah", controllers.GetWilayah)
			admin.GET("/categories", controllers.GetCategories)
			admin.GET("/commodities", controllers.GetCommodities)
			admin.POST("/categories", controllers.CreateCommodity)
			admin.POST("/categories/delete", controllers.DeleteCategory)
			admin.POST("/commodity/update", controllers.UpdateCommodity)
			admin.POST("/commodity/delete-item", controllers.DeleteCommodityItem)
		}

		// B. POTENSI LAHAN (Sesuai Log: /api/potensi-lahan)
		potensi := api.Group("/potensi-lahan")
		{
			potensi.GET("", controllers.GetPotensiLahan)
			potensi.GET("/summary", controllers.GetSummaryLahan)
			potensi.GET("/no-potential", controllers.GetNoPotentialLahan)
			potensi.GET("/filters", controllers.GetFilterOptions)
			potensi.POST("", middleware.RequireRoles(models.RoleAdmin, models.RoleOperator), controllers.CreatePotensiLahan)
		}

		// C. KELOLA LAHAN (Sesuai Log: /api/kelola-lahan)
		kelola := api.Group("/kelola-lahan")
		{
			kelola.GET("/list", controllers.GetKelolaList)
			kelola.GET("/summary", controllers.GetKelolaSummary)
			kelola.GET("/filters", controllers.GetKelolaFilterOptions)
		}

		// D. RIWAYAT LAHAN (Sesuai Log: /api/riwayat-lahan)
		riwayat := api.Group("/riwayat-lahan")
		{
			riwayat.GET("/list", controllers.GetRiwayatList)
			riwayat.GET("/summary", controllers.GetRiwayatSummary)
			riwayat.GET("/filters", controllers.GetRiwayatFilterOptions)
		}

		// E. REKAPITULASI (Sesuai Log: /api/rekapitulasi)
		rekap := api.Group("/rekapitulasi")
		{
			rekap.GET("", controllers.GetRecapData)
			rekap.GET("/export", controllers.ExportRecapExcel)
		}

		// F. SHARED VIEW
		view := api.Group("/view")
		{
			view.GET("/profile", controllers.GetProfile)
			view.GET("/jabatan", controllers.GetJabatan)
			view.GET("/tingkat", controllers.GetTingkat)
			view.GET("/wilayah", controllers.GetWilayah)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run("0.0.0.0:" + port)
}
