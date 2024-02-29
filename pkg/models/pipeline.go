package models

type Pipeline struct {
	Name    string `yaml:"name" json:"name"`
	Command string `yaml:"command" json:"command"`
}

type QSPipelineMsg struct {
	Repo     string     `json:"repo"`
	Pipeline []Pipeline `json:"pipeline"`
}
