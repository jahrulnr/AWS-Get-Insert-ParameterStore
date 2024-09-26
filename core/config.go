package core

const FilePath = "data"
const FileName = "list.json"
const FileNameGenerate = "newList.json"
const InitialParameter = "/prod/cms/"
const NewParameter = "/prod/cms/"
const ARNParameterStore = "arn:aws:ssm:ap-southeast-3::parameter"

type InsertPayload struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
	Type  string `json:"Type"`
}

type GetResponse struct {
	Parameters []InsertPayload `json:"Parameters"`
}

type TaskDefinitionEnvironment struct {
	Name      string `json:"name"`
	ValueFrom string `json:"valueFrom"`
}

type TaskDefinitionSecret struct {
	Secrets []TaskDefinitionEnvironment `json:"secrets"`
}
