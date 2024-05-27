package mysql

import (
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func InitFlink(hostname string, port int, username string, password string, dbname string) (*gorm.DB, error) {
	tmp, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, hostname, port, dbname))
	if err != nil {
		return nil, err
	}
	return tmp, nil
}
