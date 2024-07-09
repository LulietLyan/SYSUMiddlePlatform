package main

import (
	"backend/config"
	"backend/control"
	"backend/mysql"
	"backend/router"
	"log"
	"net/http"
	"os"

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
	control.InitSysInfo()
	if err := Execute(); err != nil {
		log.Println("start fail:", err.Error())
		os.Exit(-1)
	}
}

func Execute() error {
	// viper 用来取 config/dev.yaml中的参数
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		// mysql.DB.AutoMigrate(&models.User{}, &models.PresetBackground{}) // 将数据库的表自动映射为User
		if _, err := mysql.Init( //建立连接
			viper.GetString("db.hostname"), // 用viper将对应的参数取出来
			viper.GetInt("db.port"),
			viper.GetString("db.username"),
			viper.GetString("db.password"),
			viper.GetString("db.dbname"),
		); err != nil {
			return err
		}
		if err := control.InitDataSync(); err != nil {
			return err
		}

		// 最后别忘了把连接关了
		defer mysql.DB.Close()
		defer mysql.DB_Authorize.Close()
		defer mysql.DB_Demon.Close()
		defer mysql.DB_flink.Close()
		// defer mysql.SshDatabaseClient.Close()

		r := router.SetupRouter() // 初始化路由

		if err := r.Run(":2020"); err != nil {
			return err
		}
		port := viper.GetString("port")
		log.Println("port = *** =", port)

		return http.ListenAndServe(port, nil) // listen and serve
	}

	return rootCmd.Execute()
}
