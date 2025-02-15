package db_manager

import (
	"bdv-avito-merch/libs/4_common/smart_context"
	"errors"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	maxConns       = 1000
	maxConLifeTime = 120 * time.Second
	maxConIdleTime = 30 * time.Second
)

type DbManager struct {
	db        *gorm.DB
	jwtSecret string
}

func NewDbManager(sctx smart_context.ISmartContext) (*DbManager, error) {
	databaseUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return nil, errors.New("DATABASE_URL is not set")
	}

	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		return nil, errors.New("JWT_SECRET is not set")
	}

	// Создаем конфигурацию пула соединений
	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		return nil, err
	}
	config.MaxConns = int32(maxConns)
	config.MaxConnLifetime = maxConLifeTime
	config.MaxConnIdleTime = maxConIdleTime

	// Создаем пул соединений
	pool, err := pgxpool.NewWithConfig(sctx.GetContext(), config)
	if err != nil {
		return nil, err
	}

	// Преобразуем pgxpool.Pool в sql.DB через stdlib
	stdDB := stdlib.OpenDBFromPool(pool)

	// Инициализируем GORM
	db, err := gorm.Open(
		postgres.New(postgres.Config{Conn: stdDB}),
		&gorm.Config{
			PrepareStmt: true,
			ConnPool:    stdDB,
			Logger:      smart_context.NewGormLoggerWrapper(sctx),
		})
	if err != nil {
		return nil, err
	}

	return &DbManager{
		db:        db,
		jwtSecret: jwtSecret,
	}, nil
}

func (dbmanager *DbManager) GetGORM() *gorm.DB {
	return dbmanager.db.Session(&gorm.Session{NewDB: true})
}

func (dbmanager *DbManager) GetJwtSecret() string {
	return dbmanager.jwtSecret
}
