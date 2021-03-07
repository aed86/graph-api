package model

type Relation struct {
	Nodes map[int64]Node `json:"nodes"`
	Links map[int64]Link `json:"links"`
}
