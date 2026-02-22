package main

import (
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "SIKAP PRESISI Backend v2.0 Online"})
	})

	// --- PUBLIC ACCESS ---
	// Sesuai IAM: Tidak ada Signup publik. User didaftarkan oleh Admin.
	r.POST("/login", controllers.Login)

	// --- PROTECTED API AREA ---
	api := r.Group("/api")
	api.Use(middleware.RequireAuth)
	{
		// A. ADMIN ONLY RESOURCE (Role 1)
		// Fokus: User Management & Master Data
		admin := api.Group("/admin")
		admin.Use(middleware.RequireRoles(models.RoleAdmin))
		{
			// Personel Management (IAM: Admin mendaftarkan user)
			admin.POST("/users", controllers.CreateUser)
			admin.GET("/users", controllers.GetUsers)
			admin.GET("/users/:id", controllers.GetUserByID)
			admin.PUT("/users/:id", controllers.UpdateUser)
			admin.DELETE("/users/:id", controllers.DeleteUser)

			// Master Jabatan
			admin.POST("/jabatan", controllers.CreateJabatan)
			admin.GET("/jabatan", controllers.GetJabatan)
			admin.PUT("/jabatan/:id", controllers.UpdateJabatan)
			admin.DELETE("/jabatan/:id", controllers.DeleteJabatan)

			admin.GET("/tingkat", controllers.GetTingkat)

			// Master Wilayah
			admin.PUT("/wilayah/:id", controllers.UpdateWilayah)
			admin.GET("/wilayah", controllers.GetWilayah)

			// Master Komoditas
			admin.GET("/categories", controllers.GetCategories)
			admin.GET("/commodities", controllers.GetCommodities)
			admin.POST("/categories", controllers.CreateCommodity)
			admin.POST("/categories/delete", controllers.DeleteCategory)
			admin.POST("/commodity/update", controllers.UpdateCommodity)
			admin.POST("/commodity/delete-item", controllers.DeleteCommodityItem)
		}

		// B. OPERATIONAL & INPUT (Role 1 & 2)
		// Fokus: Transaksi data Lahan & Laporan
		input := api.Group("/input")
		input.Use(middleware.RequireRoles(models.RoleAdmin, models.RoleOperator))
		{
			input.POST("/lahan", controllers.CreatePotensiLahan)
			// input.PUT("/lahan/:id", controllers.UpdatePotensiLahan)
			// input.DELETE("/lahan/:id", controllers.DeletePotensiLahan)
		}

		// C. GENERAL VIEW & SHARED RESOURCE (Role 1, 2, 3)
		// Fokus: Read-only data untuk Dashboard & Mobile View
		view := api.Group("/view")
		{
			view.GET("/profile", controllers.GetProfile)
			view.GET("/jabatan", controllers.GetJabatan)
			view.GET("/tingkat", controllers.GetTingkat)
			view.GET("/wilayah", controllers.GetWilayah)
			view.GET("/categories", controllers.GetCategories)
			view.GET("/commodities", controllers.GetCommodities)

			// Lahan Resource (Read-only for all authenticated)
			view.GET("/lahan", controllers.GetPotensiLahan)
			// view.GET("/lahan/filters", controllers.GetFilterOptions)
			// view.GET("/lahan/summary", controllers.GetSummaryLahan)
			// view.GET("/lahan/no-potential", controllers.GetNoPotentialLahan)
		}
	}

	r.Run()
}
