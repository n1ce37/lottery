package main

import (
	"context"
	"strconv"

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
		return common.NewAPIGatewayResponse(common.CodeParamErr), nil
	}
	size, exists := req.QueryStringParameters["s"]
	if !exists {
		return common.NewAPIGatewayResponse(common.CodeParamErr), nil
	}
	sizeInt, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		return common.NewAPIGatewayResponse(common.CodeParamErr), nil
	}

	qi := dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("#n=:v1"),
		ProjectionExpression:   aws.String("#n,#nu,#t,balls"),
		Select:                 aws.String(dynamodb.SelectSpecificAttributes),
		ExpressionAttributeNames: map[string]*string{
			"#n":  aws.String("name"),
			"#t":  aws.String("time"),
			"#nu": aws.String("number"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String(name),
			},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int64(sizeInt),
	}
	startKey, exists := req.QueryStringParameters["sk"]
	if exists {
		qi.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"name": &dynamodb.AttributeValue{
				S: aws.String(name),
			},
			"number": &dynamodb.AttributeValue{
				N: aws.String(startKey),
			},
		}
	}

	qo, err := db.Query(&qi)
	if err != nil {
		return common.NewAPIGatewayResponse(common.CodeServerErr), nil
	}
	data := []map[string]interface{}{}
	dynamodbattribute.UnmarshalListOfMaps(qo.Items, &data)

	return common.NewAPIGatewayResponse(common.CodeOK, data), nil
}

func main() {
	db = common.NewDB()

	lambda.Start(handler)
}
