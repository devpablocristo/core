package dynamodbtable

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestConfigFromEnv(t *testing.T) {
	t.Setenv("AWS_REGION", "sa-east-1")
	t.Setenv("AUDIT_DYNAMODB_ENDPOINT", "http://localhost:4566")
	t.Setenv("AUDIT_DYNAMODB_TABLE", "audit-events")

	config := ConfigFromEnv("AUDIT_DYNAMODB")
	if config.Region != "sa-east-1" {
		t.Fatalf("unexpected region: %q", config.Region)
	}
	if config.Endpoint != "http://localhost:4566" {
		t.Fatalf("unexpected endpoint: %q", config.Endpoint)
	}
	if config.Table != "audit-events" {
		t.Fatalf("unexpected table: %q", config.Table)
	}
}

func TestPutJSONAndGetJSON(t *testing.T) {
	t.Parallel()

	api := &fakeItemAPI{}
	client := NewWithAPI(api, "audit-events")

	input := struct {
		PK   string `dynamodbav:"pk"`
		SK   string `dynamodbav:"sk"`
		Kind string `dynamodbav:"kind"`
	}{
		PK:   "tenant#acme",
		SK:   "evt#1",
		Kind: "report.ready",
	}
	if err := client.PutJSON(context.Background(), input); err != nil {
		t.Fatalf("PutJSON returned error: %v", err)
	}

	api.getItem = api.putInput.Item
	var out struct {
		PK   string `dynamodbav:"pk"`
		SK   string `dynamodbav:"sk"`
		Kind string `dynamodbav:"kind"`
	}
	found, err := client.GetJSON(context.Background(), struct {
		PK string `dynamodbav:"pk"`
		SK string `dynamodbav:"sk"`
	}{PK: "tenant#acme", SK: "evt#1"}, &out)
	if err != nil {
		t.Fatalf("GetJSON returned error: %v", err)
	}
	if !found {
		t.Fatal("expected item to be found")
	}
	if out.Kind != "report.ready" {
		t.Fatalf("unexpected out item: %#v", out)
	}
}

func TestPutItemRejectsMissingTable(t *testing.T) {
	t.Parallel()

	client := NewWithAPI(&fakeItemAPI{}, "")
	err := client.PutItem(context.Background(), map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: "x"},
	})
	if err == nil {
		t.Fatal("expected missing table error")
	}
}

type fakeItemAPI struct {
	putInput *dynamodb.PutItemInput
	getItem  map[string]types.AttributeValue
}

func (api *fakeItemAPI) PutItem(_ context.Context, input *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	api.putInput = input
	return &dynamodb.PutItemOutput{}, nil
}

func (api *fakeItemAPI) GetItem(_ context.Context, input *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	_ = input
	return &dynamodb.GetItemOutput{
		Item: api.getItem,
	}, nil
}

func TestMarshalShape(t *testing.T) {
	t.Parallel()

	item, err := attributevalue.MarshalMap(struct {
		PK string `dynamodbav:"pk"`
	}{PK: "tenant#acme"})
	if err != nil {
		t.Fatalf("MarshalMap returned error: %v", err)
	}
	if got := item["pk"].(*types.AttributeValueMemberS).Value; got != "tenant#acme" {
		t.Fatalf("unexpected pk: %q", got)
	}
}
