package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/pymq/go-test-assignment/handler"
	"github.com/spf13/viper"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"time"
)

func readConfig() {
	viper.SetConfigFile("config.toml")
	viper.SetConfigType("toml")

	viper.SetDefault("port", ":443")
	viper.SetDefault("production", true)
	viper.SetDefault("verbose", false)
	viper.SetDefault("returned_tvs_filename", "returned_tvs.xml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error while loading config file: %s \n", err)
	}
}

func main() {
	readConfig()
	db, err := sqlx.Connect("postgres", viper.GetString("db"))
	if err != nil {
		log.Fatalln(err)
	}

	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		filename := viper.GetString("returned_tvs_filename")
		processReturnedTvs(db, filename)
		for _ = range ticker.C {
			processReturnedTvs(db, filename)
		}
	}()

	h := &handler.Handler{DB: db}
	e := echo.New()

	if viper.GetBool("production") {
		e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(viper.GetStringSlice("hosts")...)
		e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
	}

	// Middleware

	if viper.GetBool("production") {
		e.Pre(middleware.HTTPSRedirect())
		e.Use(middleware.Recover())
	}
	if viper.GetBool("verbose") {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	// TV
	e.GET("/api/tvs/", h.GetTvs)
	e.POST("/api/tvs/", h.CreateTv)
	e.GET("/api/tvs/:id", h.GetTvById)
	e.PUT("/api/tvs/:id", h.PutTvById)
	e.DELETE("/api/tvs/:id", h.DeleteTvById)

	// Start
	if viper.GetBool("production") {
		go func() {
			e.Logger.Fatal(e.Start(":80"))
		}()
		e.Logger.Fatal(e.StartAutoTLS(viper.GetString("port")))
	} else {
		e.Logger.Fatal(e.Start(viper.GetString("port")))
	}
}
