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
	queryOut []map[string]types.AttributeValue
	deleted  *dynamodb.DeleteItemInput
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

func (api *fakeItemAPI) DeleteItem(_ context.Context, input *dynamodb.DeleteItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	api.deleted = input
	return &dynamodb.DeleteItemOutput{}, nil
}

func (api *fakeItemAPI) Query(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	return &dynamodb.QueryOutput{Items: api.queryOut}, nil
}

func (api *fakeItemAPI) ListTables(_ context.Context, _ *dynamodb.ListTablesInput, _ ...func(*dynamodb.Options)) (*dynamodb.ListTablesOutput, error) {
	return &dynamodb.ListTablesOutput{TableNames: []string{"audit-events"}}, nil
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

type fakePagingAPI struct {
	fakeItemAPI
	pages []dynamodb.QueryOutput
	idx   int
}

func (f *fakePagingAPI) Query(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	if f.idx >= len(f.pages) {
		return &dynamodb.QueryOutput{}, nil
	}
	out := f.pages[f.idx]
	f.idx++
	return &out, nil
}

func TestQueryJSONAll_Paginates(t *testing.T) {
	t.Parallel()

	pk := "asset-1"
	item1, err := attributevalue.MarshalMap(struct {
		AssetID string `dynamodbav:"assetId"`
		SK      string `dynamodbav:"sk"`
	}{AssetID: pk, SK: "ART#001#a"})
	if err != nil {
		t.Fatal(err)
	}
	item2, err := attributevalue.MarshalMap(struct {
		AssetID string `dynamodbav:"assetId"`
		SK      string `dynamodbav:"sk"`
	}{AssetID: pk, SK: "ART#002#b"})
	if err != nil {
		t.Fatal(err)
	}
	lek := map[string]types.AttributeValue{
		"assetId": &types.AttributeValueMemberS{Value: pk},
		"sk":      &types.AttributeValueMemberS{Value: "ART#001#a"},
	}
	api := &fakePagingAPI{
		pages: []dynamodb.QueryOutput{
			{Items: []map[string]types.AttributeValue{item1}, LastEvaluatedKey: lek},
			{Items: []map[string]types.AttributeValue{item2}},
		},
	}
	client := NewWithAPI(api, "artifacts")

	var out []struct {
		AssetID string `dynamodbav:"assetId"`
		SK      string `dynamodbav:"sk"`
	}
	if err := client.QueryJSONAll(context.Background(), &dynamodb.QueryInput{}, &out); err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 || out[0].SK != "ART#001#a" || out[1].SK != "ART#002#b" {
		t.Fatalf("%+v", out)
	}
}

func TestDeleteQueryAndHealth(t *testing.T) {
	t.Parallel()

	api := &fakeItemAPI{}
	client := NewWithAPI(api, "audit-events")

	if err := client.DeleteJSON(context.Background(), struct {
		PK string `dynamodbav:"pk"`
	}{PK: "tenant#acme"}); err != nil {
		t.Fatalf("DeleteJSON returned error: %v", err)
	}
	if api.deleted == nil {
		t.Fatal("expected delete input")
	}

	queryItems, err := attributevalue.MarshalList([]struct {
		PK string `dynamodbav:"pk"`
	}{
		{PK: "tenant#acme"},
	})
	if err != nil {
		t.Fatalf("MarshalList returned error: %v", err)
	}
	api.queryOut = []map[string]types.AttributeValue{
		queryItems[0].(*types.AttributeValueMemberM).Value,
	}
	var out []struct {
		PK string `dynamodbav:"pk"`
	}
	if err := client.QueryJSON(context.Background(), &dynamodb.QueryInput{}, &out); err != nil {
		t.Fatalf("QueryJSON returned error: %v", err)
	}
	if len(out) != 1 || out[0].PK != "tenant#acme" {
		t.Fatalf("unexpected query output: %#v", out)
	}

	if err := client.Health(context.Background()); err != nil {
		t.Fatalf("Health returned error: %v", err)
	}
}
