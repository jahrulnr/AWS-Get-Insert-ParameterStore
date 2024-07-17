package core

const FilePath = "data"
const FileName = "list.json"
const FileNameGenerate = "newList.json"
const InitialParameter = "/preprod/drupal"
const NewParameter = "/preprod/outpost/drupal"

type InsertPayload struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
	Type  string `json:"Type"`
}

type GetPayload struct {
	Names []string
}

type GetResponse struct {
	Parameters []InsertPayload `json:"Parameters"`
}
