package response

import (
	"encoding/json"

	"github.com/martini-contrib/render"
)

const SuccessStatusCode = 200

type ResponseRenderer struct {

}

type Response struct {
	Success bool   `json:"success"`
	Result  interface{} `json:"result,omitempty"`
	Error   string `json:"error,omitempty"`
}

func Result(r render.Render, result interface{}) {
	res := Response{
		Result:  result,
		Success: true,
	}

	r.JSON(200, res)
}

func Error(r render.Render, msg string, status int) {
	res, _ := json.Marshal(Response{
		Error:   msg,
		Success: false,
	})

	if status == 0 {
		status = SuccessStatusCode
	}

	r.JSON(status, res)
}
