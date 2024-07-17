package core

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var jsonResponse GetResponse

func GetParameterStore() {
	os.Mkdir(FilePath, 0755)
	executeGetParameterStore()
}

func executeGetParameterStore() {
	cmd := exec.Command("aws", "ssm", "get-parameters-by-path", "--path", InitialParameter, "--recursive", "--with-decryption")
	cmd.Env = os.Environ()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	_ = cmd.Run()
	if stdout.String() != "[]" {
		file, err := os.OpenFile(FilePath+"/"+FileName, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			log.Fatal(err)
		}
		var jsonNewResponse GetResponse
		json.Unmarshal(stdout.Bytes(), &jsonNewResponse)
		jsonResponse.Parameters = append(jsonResponse.Parameters, jsonNewResponse.Parameters...)

		joinContents, _ := json.MarshalIndent(jsonResponse, "", "  ")
		_, err = file.Write(joinContents)
		if err != nil {
			log.Println(err)
		}
		file.Close()
		time.Sleep(time.Second * time.Duration(5))
	}
}

func GenerateList() {
	list := getPayloadParameterStore(FilePath, FileName)
	newJson := GetResponse{
		Parameters: list,
	}

	file, err := os.OpenFile(FilePath+"/"+FileNameGenerate, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}

	joinContents, _ := json.MarshalIndent(newJson, "", "  ")
	_, err = file.Write(joinContents)
	if err != nil {
		log.Println(err)
	}
	file.Close()
	time.Sleep(time.Second * time.Duration(5))
}

func InsertParameterStore() {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	list := getPayloadParameterStore(FilePath, FileNameGenerate)
	for _, payload := range list {
		json, _ := json.Marshal(payload)
		log.Println("aws", "ssm", "put-parameter", "--cli-input-json", "'"+string(json)+"'", "--overwrite")
		cmd := exec.Command("aws", "ssm", "put-parameter", "--cli-input-json", string(json), "--overwrite")
		cmd.Env = os.Environ()
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			log.Println(stderr.String())
			log.Fatalln(err)
		}
	}

	log.Println(stdout.String())
}

func getPayloadParameterStore(filePath, fileName string) []InsertPayload {
	var payload GetResponse
	getParameter := getParameterStore(filePath, fileName)
	json.Unmarshal(getParameter, &payload)

	var inserPayloads []InsertPayload
	for _, paramList := range payload.Parameters {
		inserPayloads = append(inserPayloads, InsertPayload{
			Name:  strings.ReplaceAll(paramList.Name, InitialParameter, NewParameter),
			Value: paramList.Value,
			Type:  paramList.Type,
		})
	}

	return inserPayloads
}

func getParameterStore(filePath, fileName string) []byte {
	file, err := os.Open(filepath.Join(filePath, fileName))
	if err != nil {
		log.Fatalln(err)
	}

	var jsonContents []byte
	jsonContents, _ = io.ReadAll(file)

	if err != nil {
		log.Fatalln(err)
	}
	file.Close()
	return jsonContents
}
