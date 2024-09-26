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
	cmd := exec.Command("/usr/local/bin/aws", "ssm", "get-parameters-by-path", "--path", InitialParameter, "--recursive", "--with-decryption")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "AWS_REGION=ap-southeast-3")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(stderr.String())
	}
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
	// create to json
	list := getPayloadFromParameterStore(FilePath, FileName)
	newJson := GetResponse{
		Parameters: list,
	}

	file, err := os.OpenFile(FilePath+"/"+FileNameGenerate, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}

	joinContents, _ := json.MarshalIndent(newJson, "", "    ")
	_, err = file.Write(joinContents)
	if err != nil {
		log.Println(err)
	}
	file.Close()

	// create to env
	listEnv := getEnvFromParameterStore(FilePath, FileName)
	file, err = os.OpenFile(FilePath+"/"+FileNameGenerate+".env", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}

	joinContents = []byte(strings.Join(listEnv, "\n"))
	_, err = file.Write(joinContents)
	if err != nil {
		log.Println(err)
	}

	listTDF := getTDFFromParameterStore(FilePath, FileName)
	file, err = os.OpenFile(FilePath+"/"+FileNameGenerate+".tdf.json", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}

	joinContents, _ = json.MarshalIndent(listTDF, "", "    ")
	_, err = file.Write(joinContents)
	if err != nil {
		log.Println(err)
	}
}

func InsertParameterStore() {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	list := getPayloadFromParameterStore(FilePath, FileNameGenerate)
	for _, payload := range list {
		json, _ := json.Marshal(payload)
		log.Println("/usr/local/bin/aws", "ssm", "put-parameter", "--cli-input-json", "'"+string(json)+"'", "--overwrite", "--region", "ap-southeast-3")
		cmd := exec.Command("/usr/local/bin/aws", "ssm", "put-parameter", "--cli-input-json", string(json), "--overwrite", "--region", "ap-southeast-3")
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

// utils
func getPayloadFromParameterStore(filePath, fileName string) []InsertPayload {
	var payload GetResponse
	getParameter := getFromParameterStore(filePath, fileName)
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

func getEnvFromParameterStore(filePath, fileName string) []string {
	var payload GetResponse
	getParameter := getFromParameterStore(filePath, fileName)
	json.Unmarshal(getParameter, &payload)

	var inserPayloads []string
	for _, paramList := range payload.Parameters {
		inserPayloads = append(inserPayloads,
			strings.ReplaceAll(paramList.Name, InitialParameter, "")+"="+paramList.Value,
		)
	}

	return inserPayloads
}

func getTDFFromParameterStore(filePath, fileName string) TaskDefinitionSecret {
	var payload GetResponse
	getParameter := getFromParameterStore(filePath, fileName)
	json.Unmarshal(getParameter, &payload)

	var taskDef []TaskDefinitionEnvironment
	for _, paramList := range payload.Parameters {
		taskDef = append(taskDef, TaskDefinitionEnvironment{
			Name:      strings.ReplaceAll(paramList.Name, InitialParameter, ""),
			ValueFrom: ARNParameterStore + paramList.Name,
		},
		)
	}

	taskDefinition := TaskDefinitionSecret{
		Secrets: taskDef,
	}
	return taskDefinition
}

func getFromParameterStore(filePath, fileName string) []byte {
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
