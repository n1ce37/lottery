package common

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type Code int

const (
	CodeOK Code = iota
	CodeParamErr
	CodeServerErr
)

func NewAPIGatewayResponse(code Code, data ...interface{}) events.APIGatewayProxyResponse {
	bodyContent := map[string]interface{}{
		"code": code,
	}
	if code == CodeOK {
		bodyContent["data"] = data[0]
	}
	bodyBytes, _ := json.Marshal(bodyContent)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            string(bodyBytes),
		IsBase64Encoded: false,
	}
}
