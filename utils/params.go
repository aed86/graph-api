package utils

import (
	"errors"
	"strconv"

	"github.com/go-martini/martini"

	"github.com/aed86/amboss-graph-api/constants"
	"github.com/aed86/amboss-graph-api/model"
)

func GetId(params martini.Params) (int64, error) {
	v, ok := params["id"]
	if !ok || v == "" {
		return 0, errors.New("id is not provided")
	}

	val, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, errors.New("id is not provided")
	}
	return val, nil
}

func GetLimitFromReq(payload model.ReqIn) int64 {
	if payload.Limit == nil {
		return constants.NodeGetLimit
	}

	return *payload.Limit
}

func GetLimit(params martini.Params) int64 {
	v, ok := params["limit"]
	if !ok || v == "" {
		return constants.NodeGetLimit
	}

	val, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return constants.NodeGetLimit
	}
	return val
}
