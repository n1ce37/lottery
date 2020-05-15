package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/go/myService/common"
)

const tableName = "testA"

var db *dynamodb.DynamoDB

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	name, exists := req.PathParameters["name"]
	if !exists {
		log.Println("path param name is not exist")
		return common.NewAPIGatewayResponse(common.CodeParamErr), nil
	}
	number, exists := req.QueryStringParameters["n"]
	if !exists {
		log.Println("query param n is not exist")
		return common.NewAPIGatewayResponse(common.CodeParamErr), nil
	}

	gii := dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"name": &dynamodb.AttributeValue{
				S: aws.String(name),
			},
			"number": &dynamodb.AttributeValue{
				N: aws.String(number),
			},
		},
	}

	gio, err := db.GetItem(&gii)
	if err != nil {
		// TODO more detail message
		log.Println(err)

		return common.NewAPIGatewayResponse(common.CodeServerErr), nil
	}
	data := map[string]interface{}{}
	dynamodbattribute.UnmarshalMap(gio.Item, &data)

	return common.NewAPIGatewayResponse(common.CodeOK, data), nil
}

func main() {
	db = common.NewDB()

	lambda.Start(handler)
}
