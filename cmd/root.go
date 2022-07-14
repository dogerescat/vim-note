package cmd

import (
	"fmt"
	"log"
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
			// Do Stuff Here
			fmt.Println("vim-note")
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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/vim-note/config.toml)")
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		cur, err := user.Current()
		if err != nil {
			log.Fatalln(err)
		}
		viper.AddConfigPath(cur.HomeDir + "/vim-note")
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
	}
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	storage = firebase.NewStorage(config.Firebase)
}
