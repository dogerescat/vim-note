package cmd

import (
	"fmt"
	"os"
	"os/user"

	"github.com/dogerescat/vim-note/storage/firebase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "vimnote",
		Short: "vim-note",
		Long:  "vim-note root command",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Commands:")
			fmt.Println(" new                   create new memo")
			fmt.Println(" list                  show memo list")
			fmt.Println(" edit [filename]       edit memo")
		},
	}
)

type Config struct {
	Firebase firebase.Config
}

type Storage interface {
	Upload(file *os.File)
	Download(fileName string) []byte
	GetFileList() ([]string, error)
	List() error
}

var cfgFile string
var config Config
var storage Storage

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/vim-note/config.toml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		cur, err := user.Current()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		viper.AddConfigPath(cur.HomeDir + "/vim-note")
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
	}
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	storage = firebase.NewStorage(config.Firebase)
}
