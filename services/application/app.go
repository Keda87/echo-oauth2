package application

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	pg "github.com/vgarvardt/go-oauth2-pg/v4"
	"github.com/vgarvardt/go-pg-adapter/sqladapter"

	"github.com/Keda87/echo-oauth2/databases"
	"github.com/Keda87/echo-oauth2/services/api/authentications"
	"github.com/Keda87/echo-oauth2/services/api/users"
	"github.com/Keda87/echo-oauth2/services/config"
	"github.com/Keda87/echo-oauth2/services/helper"
)

type App struct {
	config       *config.Config
	DBManager    *databases.Manager
	E            *echo.Echo
	Oauth2Server *server.Server
}

func New(config *config.Config) *App {
	app := &App{
		E:         echo.New(),
		config:    config,
		DBManager: &databases.Manager{},
	}
	app.E.Validator = &helper.CustomValidator{Validator: validator.New()}

	app.initDB()
	app.initOauth2Server()
	app.initMiddleware()
	app.initRoutes()

	return app
}

func (a *App) initDB() {
	dataSource := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		a.config.DBUser,
		a.config.DBPass,
		a.config.DBHost,
		a.config.DBPort,
		a.config.DBName,
	)
	db, err := sqlx.Open("pgx", dataSource)
	if err != nil {
		a.E.Logger.Fatal(err)
	}

	// Connection pooling.
	db.SetMaxIdleConns(a.config.DBMaxIdleConnections)
	db.SetMaxOpenConns(a.config.DBMaxOpenConnections)

	a.DBManager.DB = db
}

func (a *App) initMiddleware() {
	a.E.Use(middleware.Recover())
	a.E.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAcceptEncoding},
	}))
}

func (a *App) initRoutes() {
	// Repository definitions.
	userRepo := users.NewRepository()

	// Service definitions.
	userSvc := users.NewService(a.DBManager.DB, userRepo)

	// Controller definitions.
	authController := authentications.NewController(a.Oauth2Server)
	userController := users.NewController(userSvc)

	v1 := a.E.Group("/v1")

	v1users := v1.Group("/users")
	v1users.POST("", userController.HandleRegisterUser)

	v1oauth := v1.Group("/oauth")
	v1oauth.POST("/token", authController.HandleObtainToken)
}

func (a *App) initOauth2Server() {
	manager := manage.NewDefaultManager()
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	adapter := sqladapter.NewX(a.DBManager.DB)
	clientStore, _ := pg.NewClientStore(adapter, pg.WithClientStoreLogger(a.E.Logger))
	tokenStore, _ := pg.NewTokenStore(adapter, pg.WithTokenStoreGCInterval(time.Minute))

	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(false)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	a.Oauth2Server = srv
}

func (a *App) StartServer() {
	a.E.HideBanner = true

	go func() {
		if err := a.E.Start(":" + a.config.APPPort); err != nil {
			a.E.Logger.Info("shutting down the server")
		}
	}()

	// Gracefully shutdown the server.
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGABRT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.E.Shutdown(ctx); err != nil {
		a.E.Logger.Fatal(err)
	}
}

func (a *App) PreStopServer() {
	fmt.Println("Closing connections...")
	a.DBManager.DB.Close()
}
