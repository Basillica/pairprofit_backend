package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/gin-gonic/gin"
)

//https://dynobase.dev/dynamodb-golang-query-examples/#setup
// https://docs.localstack.cloud/get-started/#docker

func CreateTable(c *gin.Context) {
	dbClient := c.MustGet("dbClient").(*dynamodb.Client)
	out, err := dbClient.CreateTable(c, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName:   aws.String("my-table"),
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(out)
}

func GetAll(c *gin.Context) {
	dbClient := c.MustGet("dbClient").(*dynamodb.Client)
	out, err := dbClient.Scan(c, &dynamodb.ScanInput{
		TableName: aws.String("my-table"),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(out.Items)
}

func FilterItem(c *gin.Context) {
	dbClient := c.MustGet("dbClient").(*dynamodb.Client)
	out, err := dbClient.Scan(c, &dynamodb.ScanInput{
		TableName:        aws.String("my-table"),
		FilterExpression: aws.String("attribute_not_exists(deletedAt) AND contains(firstName, :firstName)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":firstName": &types.AttributeValueMemberS{Value: "John"},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(out.Items)
}

func GetItem(c *gin.Context) {
	dbClient := c.MustGet("dbClient").(*dynamodb.Client)
	out, err := dbClient.GetItem(c, &dynamodb.GetItemInput{
		TableName: aws.String("my-table"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: "123"},
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(out.Item)
}

func CreateItemFromStruct(c *gin.Context) {
	dbClient := c.MustGet("dbClient").(*dynamodb.Client)
	key := struct {
		ID string `dynamodbav:"id" json:"id"`
	}{ID: "123"}
	avs, err := attributevalue.MarshalMap(key)
	if err != nil {
		panic(err)
	}

	out, err := dbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("my-table"),
		Key:       avs,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(out.Item)
}

func BatchGetItem(c *gin.Context) {
	dbClient := c.MustGet("dbClient").(*dynamodb.Client)
	out, err := dbClient.BatchGetItem(context.TODO(), &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			"my-table": {
				Keys: []map[string]types.AttributeValue{
					{
						"id": &types.AttributeValueMemberS{Value: "123"},
					},
					{
						"id": &types.AttributeValueMemberS{Value: "123"},
					},
				},
			},
			"other-table": {
				Keys: []map[string]types.AttributeValue{
					{
						"id": &types.AttributeValueMemberS{Value: "abc"},
					},
					{
						"id": &types.AttributeValueMemberS{Value: "abd"},
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(out.Responses)
}

func PutItem(c *gin.Context) {
	dbClient := c.MustGet("dbClient").(*dynamodb.Client)
	out, err := dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("my-table"),
		Item: map[string]types.AttributeValue{
			"id":    &types.AttributeValueMemberS{Value: "12346"},
			"name":  &types.AttributeValueMemberS{Value: "John Doe"},
			"email": &types.AttributeValueMemberS{Value: "john@doe.io"},
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(out.Attributes)
}

func BatchPutItem(c *gin.Context) {
	dbClient := c.MustGet("dbClient").(*dynamodb.Client)
	out, err := dbClient.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			"TableOne": {
				{
					DeleteRequest: &types.DeleteRequest{
						Key: map[string]types.AttributeValue{
							"id": &types.AttributeValueMemberS{Value: "123"},
						},
					},
				},
				{
					PutRequest: &types.PutRequest{
						Item: map[string]types.AttributeValue{
							"id":    &types.AttributeValueMemberS{Value: "234"},
							"name":  &types.AttributeValueMemberS{Value: "dynamobase"},
							"email": &types.AttributeValueMemberS{Value: "dynobase@dynobase.dev"},
						},
					},
				},
			},
			"TableTwo": {
				{
					PutRequest: &types.PutRequest{
						Item: map[string]types.AttributeValue{
							"id":    &types.AttributeValueMemberS{Value: "456"},
							"name":  &types.AttributeValueMemberS{Value: "dynamobase"},
							"email": &types.AttributeValueMemberS{Value: "dynobase@dynobase.dev"},
						},
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(out)
}

func GetItems(c *gin.Context) {
	dbClient := c.MustGet("dbClient").(*dynamodb.Client)
	out, err := dbClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String("my-table"),
		KeyConditionExpression: aws.String("id = :hashKey and #date > :rangeKey"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":hashKey":  &types.AttributeValueMemberS{Value: "123"},
			":rangeKey": &types.AttributeValueMemberN{Value: "20150101"},
		},
		ExpressionAttributeNames: map[string]string{
			"#date": "date",
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(out.Items)
}

func QueryAnIndex(c *gin.Context) {
	dbClient := c.MustGet("dbClient").(*dynamodb.Client)
	out, err := dbClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String("my-table"),
		IndexName:              aws.String("GSI1"),
		KeyConditionExpression: aws.String("gsi1pk = :gsi1pk and gsi1sk > :gsi1sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":gsi1pk": &types.AttributeValueMemberS{Value: "123"},
			":gsi1sk": &types.AttributeValueMemberN{Value: "20150101"},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(out.Items)
}

func QueryWithSorting(c *gin.Context) {
	dbClient := c.MustGet("dbClient").(*dynamodb.Client)
	out, err := dbClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String("my-table"),
		IndexName:              aws.String("Index"),
		KeyConditionExpression: aws.String("id = :hashKey and #date > :rangeKey"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":hashKey":  &types.AttributeValueMemberS{Value: "123"},
			":rangeKey": &types.AttributeValueMemberN{Value: "20150101"},
		},
		ExpressionAttributeNames: map[string]string{
			"#date": "date",
		},
		ScanIndexForward: aws.Bool(true), // true or false to sort by "date" Sort/Range key ascending or descending
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(out.Items)
}

func QueryAndScanWithPagination(c *gin.Context) {
	dbClient := c.MustGet("dbClient").(*dynamodb.Client)
	p := dynamodb.NewQueryPaginator(dbClient, &dynamodb.QueryInput{
		TableName:              aws.String("my-table"),
		Limit:                  aws.Int32(1),
		KeyConditionExpression: aws.String("id = :hashKey and #date > :rangeKey"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":hashKey":  &types.AttributeValueMemberS{Value: "123"},
			":rangeKey": &types.AttributeValueMemberN{Value: "20150101"},
		},
		ExpressionAttributeNames: map[string]string{
			"#date": "date",
		},
	})

	type Item struct {
		ID   string `dynamodbav:"id"`
		Date int    `dynamodbav:"date"`
	}
	var items []Item
	for p.HasMorePages() {
		out, err := p.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}

		var pItems []Item
		err = attributevalue.UnmarshalListOfMaps(out.Items, &pItems)
		if err != nil {
			panic(err)
		}

		items = append(items, pItems...)
	}

	fmt.Println(items)
}

func UpdateItem(c *gin.Context) {
	dbClient := c.MustGet("dbClient").(*dynamodb.Client)
	out, err := dbClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String("my-table"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: "123"},
		},
		UpdateExpression: aws.String("set firstName = :firstName"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":firstName": &types.AttributeValueMemberS{Value: "John McNewname"},
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(out.Attributes)
}

func ConditionalUpdate(c *gin.Context) {
	dbClient := c.MustGet("dbClient").(*dynamodb.Client)
	out, err := dbClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String("my-table"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: "123"},
		},
		UpdateExpression: aws.String("set firstName = :firstName"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":firstName": &types.AttributeValueMemberS{Value: "John McNewname"},
			":company":   &types.AttributeValueMemberS{Value: "Apple"},
		},
		ConditionExpression: aws.String("attribute_not_exists(deletedAt) and company = :company"),
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(out.Attributes)
}

func UpdateWithExpressionBuilder(c *gin.Context) {
	dbClient := c.MustGet("dbClient").(*dynamodb.Client)
	expr, err := expression.NewBuilder().WithUpdate(
		expression.Set(
			expression.Name("firstName"),
			expression.Value("John McNewname"),
		),
	).WithCondition(
		expression.And(
			expression.AttributeNotExists(
				expression.Name("deletedAt"),
			),
			expression.Equal(
				expression.Name("company"),
				expression.Value("Apple"),
			),
		),
	).Build()
	if err != nil {
		panic(err)
	}

	out, err := dbClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String("my-table"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: "123"},
		},
		UpdateExpression: expr.Update(),
		// ExpressionAttributeNames:  expr.Names(),
		// ExpressionAttributeValues: expr.Values(),
		ConditionExpression: expr.Condition(),
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(out.Attributes)
}

func IncrementItemAttr(c *gin.Context) {
	dbClient := c.MustGet("dbClient").(*dynamodb.Client)
	out, err := dbClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String("my-table"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: "123"},
		},
		UpdateExpression: aws.String("set score = score + :value"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":value": &types.AttributeValueMemberN{Value: "1"},
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(out.Attributes)
}

func DeleteItem(c *gin.Context) {
	dbClient := c.MustGet("dbClient").(*dynamodb.Client)
	out, err := dbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String("my-table"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: "123"},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(out.Attributes)
}

func DeleteAllItem(c *gin.Context) {
	dbClient := c.MustGet("dbClient").(*dynamodb.Client)
	p := dynamodb.NewScanPaginator(dbClient, &dynamodb.ScanInput{
		TableName: aws.String("my-table"),
	})

	for p.HasMorePages() {
		out, err := p.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}

		for _, item := range out.Items {
			_, err = dbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
				TableName: aws.String("my-table"),
				Key: map[string]types.AttributeValue{
					"pk": item["pk"],
					"sk": item["sk"],
				},
			})
			if err != nil {
				panic(err)
			}
		}
	}
}
