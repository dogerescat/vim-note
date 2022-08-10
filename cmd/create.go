/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "new",
	Short: "create new file",
	Long:  "create new file",
	Run: func(cmd *cobra.Command, args []string) {
		fileList, _ := storage.GetFileList()
		fileName := args[0]
		if exist(fileName, fileList) {
			fmt.Println("already exist this file")
			os.Exit(0)
		}
		c := exec.Command("vim", fileName)
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		c.Stdin = os.Stdin
		err := c.Run()
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
	rootCmd.AddCommand(createCmd)
}
