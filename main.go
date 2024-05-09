package main

import (
	"backend/config"
	"backend/models"
	"backend/mysql"
	"backend/router"
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{}
)

func initConfig() {
	config.MustInit(os.Stdout, cfgFile) // 配置初始化
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config/dev.yaml", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().Bool("debug", true, "开启debug")
	viper.SetDefault("gin.mode", rootCmd.PersistentFlags().Lookup("debug"))
}

func main() {
	if err := Execute(); err != nil {
		log.Println("start fail:", err.Error())
		os.Exit(-1)
	}
}

func Execute() error {
	// viper 用来取 config/dev.yaml中的参数
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		// _, err := mysql.DialWithPassword(
		// 	viper.GetString("ssh.Host"),
		// 	viper.GetInt("ssh.Port"),
		// 	viper.GetString("ssh.User"),
		// 	viper.GetString("ssh.Password"),
		// )

		// if err != nil {
		// 	return err
		// }

		// _, err = mysql.Init(
		// 	viper.GetString("db.hostname"),
		// 	viper.GetInt("db.port"),
		// 	viper.GetString("db.username"),
		// 	viper.GetString("db.password"),
		// 	viper.GetString("db.dbname"),
		// ) // 用viper将对应的参数取出来
		// if err != nil {
		// 	return err
		// }

		// mysql.DB.AutoMigrate(&models.User{}, &models.PresetBackground{}) // 将数据库的表自动映射为User
		var err error
		mysql.DB, err = gorm.Open("mysql", "root:654321@tcp(47.121.29.57:3307)/mydb?charset=utf8mb4&parseTime=True&loc=Local")
		if err != nil {
			return err
		}
		mysql.DB.Set("gorm:table_options", "charset=utf8mb4").
			AutoMigrate(&models.User{}).
			AutoMigrate(&models.PresetBackground{})

		// 最后别忘了把连接关了
		defer mysql.DB.Close()
		// defer mysql.SshDatabaseClient.Close()

		r := router.SetupRouter() // 初始化路由
		err = r.Run(":2020")
		if err != nil {
			return err
		}

		port := viper.GetString("port")
		log.Println("port = *** =", port)
		return http.ListenAndServe(port, nil) // listen and serve
	}

	return rootCmd.Execute()
}
