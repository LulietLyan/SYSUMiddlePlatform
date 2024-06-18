package mysql

import (
	"backend/models"
	"context"
	_ "database/sql"
	"fmt"
	"net"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/ssh"
)

var (
	DB                *gorm.DB
	DB_flink          *gorm.DB
	DB_Authorize      *gorm.DB
	SshDatabaseClient *ssh.Client
)

type ViaSSHDialer struct {
	Client *ssh.Client
}

func (v *ViaSSHDialer) Dial(context context.Context, addr string) (net.Conn, error) {
	return v.Client.Dial("tcp", addr)
}

func DialWithPassword(hostname string, port int, username string, password string) (*ssh.Client, error) {
	address := fmt.Sprintf("%s:%d", hostname, port)
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshDatabaseClient, err := ssh.Dial("tcp", address, config)

	if err != nil {
		return sshDatabaseClient, err
	}
	mysql.RegisterDialContext("mysql+tcp", (&ViaSSHDialer{Client: sshDatabaseClient}).Dial)
	SshDatabaseClient = sshDatabaseClient
	return sshDatabaseClient, nil
}

func Init(hostname string, port int, username string, password string, dbname string) (*gorm.DB, error) {
	var err error
	DB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, hostname, port, dbname))
	if err != nil {
		return nil, err
	}

	DB_Authorize, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, hostname, port, "mysql"))
	if err != nil {
		return nil, err
	}

	DB_flink, _ = InitFlink(hostname, port, username, password, "flink_target")
	return DB, DB.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(&models.User{}).
		AutoMigrate(&models.AnalyticalUser{}).
		AutoMigrate(&models.ProjectTable{}).
		AutoMigrate(&models.ProjectUser{}).
		AutoMigrate(&models.Permission{}).
		AutoMigrate(&models.DingdingProjectUser{}).
		AutoMigrate(&models.DingdingAnalyticalUser{}).
		AutoMigrate(&models.Notifications{}).
		AutoMigrate(&models.PermissionRequest{}).
		AutoMigrate(&models.Api{}).
		AutoMigrate(&models.ActivationCode{}).
		AutoMigrate(&models.ProjectMember{}).Error
}
