package model

type Relation struct {
	Nodes map[int64]Node `json:"nodes"`
	Links []Link `json:"links"`
}
