package main

import (
	"aws-param-store/core"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args != nil {
		switch os.Args[1] {
		case "getlist":
			os.Mkdir(core.FilePath, 0755)
			os.Remove(core.FilePath + "/" + core.FileName)
			core.GetParameterStore()
		case "generatelist":
			os.Remove(core.FilePath + "/" + core.FileNameGenerate)
			os.Remove(core.FilePath + "/" + core.FileNameGenerate + ".env")
			core.GenerateList()
		case "insertlist":
			core.InsertParameterStore()
		case "createFromEnv":
			core.CreatePayloadFromEnv()
		default:
			fmt.Println("Command not found")
		}

		return
	}

	fmt.Println("Available command: getlist, generatelist, createFromEnv, insertlist")
}
