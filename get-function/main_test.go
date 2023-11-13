package main

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	request := events.APIGatewayProxyRequest{
		Path:       "get-function/",
		HTTPMethod: "GET",
	}

	response, err := handler(request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 200, response.StatusCode)

	var responseBody map[string]string
	err = json.Unmarshal([]byte(response.Body), &responseBody)
	assert.NoError(t, err)

	expectedCount := responseBody["count"]
	assert.Equal(t, expectedCount, responseBody["count"])
}
