package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/dogerescat/vim-note/ui"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit file",
	Long:  "edit file",
	Run: func(cmd *cobra.Command, args []string) {
		fileList, _ := storage.GetFileList()
		var fileName string
		data := []byte("")
		if len(args) != 0 {
			if !exist(args[0], fileList) {
				fmt.Println("not exist file")
				os.Exit(1)
			}
			fileName = args[0]
		} else {
			fileName = ui.Run(fileList)
			if fileName == "" {
				os.Exit(0)
			}
		}
		data = storage.Download(fileName)
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
				fmt.Println(err.Error())
			}
			err = content.Close()
			if err != nil {
				fmt.Println(err.Error())
			}
			os.Exit(0)
		}()
		storage.Upload(content)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
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
