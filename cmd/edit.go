/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit file",
	Long:  "edit file",
	Run: func(cmd *cobra.Command, args []string) {
		fileList, _ := storage.GetFileList()
		fileName := args[0]
		data := []byte("")
		if len(args) != 0 {
			if !exist(fileName, fileList) {
				fmt.Println("not exist file")
				os.Exit(1)
			}
			data = storage.Download(args[0])
		} else {
			fmt.Println("non title")
			os.Exit(0)
		}
		err := ioutil.WriteFile(fileName, data, 0664)
		if err != nil {
			fmt.Printf(err.Error())
		}
		c := exec.Command("vim", fileName)
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		c.Stdin = os.Stdin
		err = c.Run()
		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
		}
		content, _ := os.Open(fileName)
		defer func() {
			err = os.Remove(fileName)
			if err != nil {
				fmt.Println(err)
			}
			err = content.Close()
			if err != nil {
				fmt.Println(err)
			}
			os.Exit(0)
		}()
		storage.Upload(content)

	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// editCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// editCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func exist(fileName string, fileList []string) bool {
	ok := false
	for i := 0; i < len(fileList); i++ {
		if fileList[i] == fileName {
			ok = true
		}
	}
	return ok
}
