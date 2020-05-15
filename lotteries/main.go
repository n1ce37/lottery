package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/go/myService/common"
)

const tableName = "lotteries"

var db *dynamodb.DynamoDB

func handler(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	si := dynamodb.ScanInput{
		TableName: aws.String(tableName),
		// ProjectionExpression: aws.String("#n,jackpot,nextTime,history[0].#t,history[0].balls"),
		// ExpressionAttributeNames: map[string]*string{
		// 	"#n": aws.String("name"),
		// 	"#t": aws.String("time"),
		// },
		// Select: aws.String(dynamodb.SelectSpecificAttributes),
	}
	so, err := db.Scan(&si)
	if err != nil {
		return common.NewAPIGatewayResponse(common.CodeServerErr), nil
	}

	data := []map[string]interface{}{}
	dynamodbattribute.UnmarshalListOfMaps(so.Items, &data)

	return common.NewAPIGatewayResponse(common.CodeOK, data), nil
}

func main() {
	db = common.NewDB()

	lambda.Start(handler)
}
