package main

import (
	"aws-param-store/core"
	"os"
)

func main() {
	if os.Args != nil {
		switch os.Args[1] {
		case "getlist":
			os.Mkdir(core.FilePath, 0755)
			os.Remove(core.FilePath + "/" + core.FileName)
			core.GetParameterStore()
		case "generatelist":
			os.Remove(core.FilePath + "/" + core.FileNameGenerate)
			core.GenerateList()
		case "insertlist":
			core.InsertParameterStore()
		}
	}
}
