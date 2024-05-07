package mysql

import (
	"context"
	_ "database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/ssh"
	"net"
)

var (
	DB                *gorm.DB
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
	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@mysql+tcp(%s:%d)/%s?charset=utf8&parseTime=True",
			username, password, hostname, port, dbname))
	if err != nil {
		return nil, err
	}
	DB = db
	return db, nil
}
