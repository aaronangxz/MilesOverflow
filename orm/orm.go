package orm

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	defaultLog "log"
	"os"
	"time"
)

var (
	db        *gorm.DB
	newLogger = logger.New(
		defaultLog.New(os.Stdout, "\r\n", defaultLog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
)

func DbInstance() *gorm.DB {
	if db == nil {
		ConnectMySQL()
	}
	return db
}

func ConnectMySQL() {
	URL := fmt.Sprintf("%v:%v@tcp(%v)/%v", "root", "Xuanze94", "localhost:3306", "milestracker_db")

	log.Infof("Connecting to %v", URL)
	openDb, err := gorm.Open(mysql.Open(URL), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Errorf("Error while establishing Live DB Connection: %v", err)
		panic("Failed to connect to live database!")
	}
	log.Info("Live Database connection established")
	db = openDb
}
