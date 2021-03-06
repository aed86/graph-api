package model

type PathOut struct {
	Idx            int64       `json:"idx"`
	SourceNodeName string      `json:"sourceNodeName"`
	TargetNodeName string      `json:"targetNodeName"`
	TotalCost      float64       `json:"totalCost"`
	Path           interface{} `json:"path"`
	PathCosts      interface{} `json:"path_costs"`
}
