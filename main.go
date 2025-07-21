package main

import (
	"database/sql"
	"errors"
	"fmt"
	"gold-savings/admin"
	"gold-savings/admin/components"
	"gold-savings/api"
	"gold-savings/api/middleware"
	db "gold-savings/db/sqlc"
	"gold-savings/internal/auth"
	"gold-savings/internal/config"
	p "gold-savings/internal/db"
	"gold-savings/internal/services"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	c, err := config.LoadConfig(".")
	if err != nil {
		panic(fmt.Sprintf("Could not load config: %v", err))
	}

	dbConn, err := sql.Open(c.DBDriver, p.GetDBSource(c, c.DBName))
	if err != nil {
		panic(fmt.Sprintf("Could not load DB: %v", err))
	}

	m, err := migrate.New(
		"file://db/migrations",
		p.GetDBSource(c, c.DBName),
	)
	if err != nil {
		log.Fatalf("Unable to instantiate the database schema migrator - %v", err)
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Unable to migrate up to the latest database schema - %v", err)
		}
	}

	defer func(dbConn *sql.DB) {
		err := dbConn.Close()
		if err != nil {
			panic(fmt.Sprintf("Could not close DB: %v", err))
		}
	}(dbConn)

	queries := db.New(dbConn)

	authService := auth.NewAuthService(queries, c)
	adminService := services.NewAdminService(queries,c)
	userService := services.NewUserService(queries, c)

	// Create a Gin router
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	// Serve static files
	router.Static("/static", "./static")

	router.GET("/.well-known/appspecific/com.chrome.devtools.json", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"}) // or just return 204 No Content
	})

	router.GET("/", func(ctx *gin.Context) {
		err := components.LandingPage().Render(ctx, ctx.Writer)
		if err != nil {
			ctx.String(404, "landing page component not found")
			return
		}
	})

	// Initialize API routes
	api.SetupRoutes(router.Group("/api"), authService, queries, c, userService)

	// Initialize Admin routes
	admin.SetupRoutes(router.Group("/admin"), authService, queries, adminService)

	port := c.ServerPort
	log.Printf("Server running on port %v", port)
	log.Fatal(router.Run(fmt.Sprintf(":%v", port)))
}
