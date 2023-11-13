package main

import (
	"errors"
	"log"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type MockDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
	UpdateItemOutput dynamodb.UpdateItemOutput
	UpdateItemError  error
}

func (m *MockDynamoDBClient) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	return &m.UpdateItemOutput, m.UpdateItemError
}

func t_handler(request events.APIGatewayProxyRequest, svc dynamodbiface.DynamoDBAPI) (events.APIGatewayProxyResponse, error) {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("cloud-resume-challenge"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String("visitors"),
			},
		},
		UpdateExpression: aws.String("ADD visitors :inc"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":inc": {
				N: aws.String("1"),
			},
		},
	}

	_, err := svc.UpdateItem(input)

	if err != nil {
		log.Printf("Got error calling UpdateItem: %s", err)
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{}, nil
}

func TestHandler(t *testing.T) {
	tests := []struct {
		name             string
		dynamoDBClient   dynamodbiface.DynamoDBAPI
		expectedResponse events.APIGatewayProxyResponse
		expectedError    error
	}{
		{
			name: "Successfully update item",
			dynamoDBClient: &MockDynamoDBClient{
				UpdateItemOutput: dynamodb.UpdateItemOutput{},
				UpdateItemError:  nil,
			},
			expectedResponse: events.APIGatewayProxyResponse{},
			expectedError:    nil,
		},
		{
			name: "Error updating item",
			dynamoDBClient: &MockDynamoDBClient{
				UpdateItemOutput: dynamodb.UpdateItemOutput{},
				UpdateItemError:  errors.New("DynamoDB update error"),
			},
			expectedResponse: events.APIGatewayProxyResponse{},
			expectedError:    errors.New("DynamoDB update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := t_handler(events.APIGatewayProxyRequest{}, tt.dynamoDBClient)

			if !reflect.DeepEqual(response, tt.expectedResponse) {
				t.Errorf("Expected response %v, but got %v", tt.expectedResponse, response)
			}

			if err != nil && (tt.expectedError == nil || err.Error() != tt.expectedError.Error()) {
				t.Errorf("Expected error %v, but got %v", tt.expectedError, err)
			}
		})

	}
}
