package main

import (
	"fmt"
	"go-test/internals/core/services"
	"go-test/internals/handlers"
	"go-test/internals/repositories"
	"go-test/internals/server"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

// main is the entry point of the application.
//
// It initializes the database, creates the necessary repositories,
// services, and handlers, and starts the HTTP server.
func main() {
	initConfig()
	db := initDatabase()

	postRepository := repositories.NewPostRepository(db)

	postService := services.NewPostService(postRepository)

	postHandler := handlers.NewPostHandler(postService)

	httpServer := server.NewServer(postHandler)

	httpServer.Initialize()

}

func initDatabase() *sqlx.DB {
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable&timezone=Asia/Bangkok",
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetInt("db.port"),
		viper.GetString("db.database"),
	)
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("..")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
