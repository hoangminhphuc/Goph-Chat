package db

import (
	"errors"
	"flag"
	"strings"
	"sync"
	"time"

	"github.com/hoangminhphuc/goph-chat/common/logger"
	"github.com/hoangminhphuc/goph-chat/plugin/gorm/gormdiaclects"
	"gorm.io/gorm"
)


const (
	GORM_DB_TYPE_MYSQL = "mysql"
	// Can have multiple type of db here: postgres, sqlite, mssql
	GORM_DB_TYPE_NOT_SUPPORTED = ""
)

type GormDB struct {
	name 				string
	URI 				string
	Type 				string
	db 					*gorm.DB
	isRunning 	bool
	mt 					sync.Mutex
	logger 			logger.ZapLogger
}

func NewGormDB(name string) *GormDB {
	return &GormDB{
		name: name,
		Type: name,
		isRunning: false,
		logger: logger.NewZapLogger(),
	}
}

func (gdb *GormDB) Name() string {
	return gdb.name
}


func (gdb * GormDB) InitFlags() {
	prefix := gdb.name + "-"
	gdb.logger.Log.Info("Initializing flags...")


	flag.StringVar(&gdb.URI, prefix + "gorm-db-uri", "", "URI for GORM database connection")
	flag.StringVar(&gdb.Type, prefix + "gorm-db-type", "", "Database type for GORM")
}

func (gdb *GormDB) Run() error { 
	gdb.mt.Lock()
	defer gdb.mt.Unlock()

	if !gdb.isRunning {
		dbType := GetDBType(gdb.Type)
		if dbType == GORM_DB_TYPE_NOT_SUPPORTED {
			return errors.New("gorm database type is not supported")
		}

		gdb.logger.Log.Info("Starting database...")
		
		var err error
		gdb.db, err = gdb.GetDBConnection(dbType)
		
		if err != nil {
			gdb.logger.Log.Error("Cannot connect to db through ", gdb.URI, ". ", err.Error())
			return err
		}
		gdb.isRunning = true
	}

	return nil
}

func GetDBType(dbType string) string {
	switch strings.ToLower(dbType) {
	case "mysql":
		return GORM_DB_TYPE_MYSQL
	// case "postgres":
	// 	return GORM_DB_TYPE_POSTGRE
	}

	return GORM_DB_TYPE_NOT_SUPPORTED
}


	
func (gdb *GormDB) GetDBConnection(dbType string) (*gorm.DB, error) {
	switch dbType {
	case GORM_DB_TYPE_MYSQL:
		return gormdiaclects.MySQLConnection(gdb.URI)
	// case GORM_DB_TYPE_POSTGRE:
	}
	return nil, nil
}


// New session of db connection
func (gdb *GormDB) Get() interface{} {
	newSessionDB := gdb.db.Session(&gorm.Session{NewDB: true})

	if db, err := newSessionDB.DB(); err == nil {
			db.SetMaxOpenConns(100)
			db.SetMaxIdleConns(100)
			db.SetConnMaxIdleTime(time.Hour)
	}

	return gdb.db.Debug()
}

func (gdb *GormDB) Stop() <-chan error {
	c := make(chan error, 1)

	go func() {
		gdb.mt.Lock()
		defer gdb.mt.Unlock()

		if !gdb.isRunning {
			c <- nil
			return
		}

		db, err := gdb.db.DB()
		if err == nil {
			err = db.Close()
		}

		if err == nil {
			gdb.isRunning = false
			gdb.logger.Log.Info("Database stopped.")
		}

		c <- err
	}()

	return c
}
