package app

import (
	"context"
	"net/http"
	"time"

	"retrofit/backend/internal/config"
	"retrofit/backend/internal/database"
	transport "retrofit/backend/internal/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	cfg    config.Config
	db     *gorm.DB
	redis  *redis.Client
	logger *zap.Logger
	router *gin.Engine
}

func New(cfg config.Config) (*App, error) {
	logger, err := zap.NewProduction()
	if cfg.AppEnv == "local" {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		return nil, err
	}

	db, err := database.NewMySQL(cfg.MySQLDSN)
	if err != nil {
		return nil, err
	}

	rdb := database.NewRedis(cfg.Redis)
	if err := pingRedis(rdb); err != nil {
		logger.Warn("redis unavailable", zap.Error(err))
	}

	router := transport.NewRouter(transport.RouterDeps{
		Config: cfg,
		DB:     db,
		Redis:  rdb,
		Logger: logger,
	})

	return &App{
		cfg:    cfg,
		db:     db,
		redis:  rdb,
		logger: logger,
		router: router,
	}, nil
}

func (a *App) Run() error {
	server := &http.Server{
		Addr:              a.cfg.HTTPAddr,
		Handler:           a.router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	a.logger.Info("api listening", zap.String("addr", a.cfg.HTTPAddr))
	return server.ListenAndServe()
}

func pingRedis(rdb *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return rdb.Ping(ctx).Err()
}
