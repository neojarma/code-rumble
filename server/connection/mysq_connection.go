package connection

import (
	"errors"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetConnectionMySQL() (*gorm.DB, error) {
	// DB_URL := os.Getenv("DB_URL")
	DB_URL := "root:secretrootpassword@tcp(localhost:3306)/code_rumble?charset=utf8mb4&parseTime=True&loc=Local"

	maxTries := 10
	for i := 0; i < maxTries; i++ {
		db, err := gorm.Open(mysql.Open(DB_URL), &gorm.Config{})
		if err == nil {
			log.Println("success connect to mysql")
			return db, nil
		}

		log.Println("failed to connect to mysql, try again in 1 minute")
		time.Sleep(1 * time.Minute)
	}

	return nil, errors.New("mysql connection failed after 10 minute")

}
