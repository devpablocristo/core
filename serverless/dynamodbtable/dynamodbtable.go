package dynamodbtable

import (
	"context"
	"fmt"
	"os"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Config representa la configuración mínima para una tabla DynamoDB.
type Config struct {
	Region   string
	Endpoint string
	Table    string
}

type itemAPI interface {
	PutItem(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

// Client envuelve acceso reusable a una tabla fija.
type Client struct {
	api   itemAPI
	table string
}

// ConfigFromEnv carga configuración desde env.
func ConfigFromEnv(prefix string) Config {
	prefix = normalizeEnvPrefix(prefix)
	if prefix == "" {
		return Config{
			Region:   firstNonEmpty(os.Getenv("AWS_REGION"), "us-east-1"),
			Endpoint: firstNonEmpty(os.Getenv("DYNAMODB_ENDPOINT")),
			Table:    firstNonEmpty(os.Getenv("DYNAMODB_TABLE")),
		}
	}
	return Config{
		Region:   firstNonEmpty(os.Getenv(prefix+"AWS_REGION"), os.Getenv("AWS_REGION"), "us-east-1"),
		Endpoint: firstNonEmpty(os.Getenv(prefix+"ENDPOINT"), os.Getenv("DYNAMODB_ENDPOINT")),
		Table:    firstNonEmpty(os.Getenv(prefix+"TABLE"), os.Getenv("DYNAMODB_TABLE")),
	}
}

// New crea un cliente DynamoDB desde configuración local.
func New(ctx context.Context, config Config) (*Client, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(config.Region))
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}
	return NewFromAWSConfig(awsCfg, config), nil
}

// NewFromAWSConfig crea un cliente reutilizando `aws.Config`.
func NewFromAWSConfig(awsCfg awssdk.Config, config Config) *Client {
	api := dynamodb.NewFromConfig(awsCfg, func(options *dynamodb.Options) {
		if strings.TrimSpace(config.Endpoint) != "" {
			options.BaseEndpoint = awssdk.String(strings.TrimSpace(config.Endpoint))
		}
	})
	return NewWithAPI(api, config.Table)
}

// NewWithAPI crea un cliente desde un adapter inyectado.
func NewWithAPI(api itemAPI, table string) *Client {
	return &Client{
		api:   api,
		table: strings.TrimSpace(table),
	}
}

// Table devuelve el nombre fijo de la tabla.
func (c *Client) Table() string {
	if c == nil {
		return ""
	}
	return c.table
}

// PutItem guarda un item raw.
func (c *Client) PutItem(ctx context.Context, item map[string]types.AttributeValue) error {
	if c == nil || c.api == nil {
		return fmt.Errorf("dynamodb client is nil")
	}
	if strings.TrimSpace(c.table) == "" {
		return fmt.Errorf("dynamodb table is required")
	}
	if len(item) == 0 {
		return fmt.Errorf("dynamodb item is required")
	}
	_, err := c.api.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: awssdk.String(c.table),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("put dynamodb item: %w", err)
	}
	return nil
}

// PutJSON serializa una estructura a AttributeValue map y la persiste.
func (c *Client) PutJSON(ctx context.Context, value any) error {
	item, err := attributevalue.MarshalMap(value)
	if err != nil {
		return fmt.Errorf("marshal dynamodb item: %w", err)
	}
	return c.PutItem(ctx, item)
}

// GetJSON resuelve una key tipada y deserializa el item si existe.
func (c *Client) GetJSON(ctx context.Context, key any, out any) (bool, error) {
	if c == nil || c.api == nil {
		return false, fmt.Errorf("dynamodb client is nil")
	}
	if strings.TrimSpace(c.table) == "" {
		return false, fmt.Errorf("dynamodb table is required")
	}

	marshaledKey, err := attributevalue.MarshalMap(key)
	if err != nil {
		return false, fmt.Errorf("marshal dynamodb key: %w", err)
	}
	output, err := c.api.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: awssdk.String(c.table),
		Key:       marshaledKey,
	})
	if err != nil {
		return false, fmt.Errorf("get dynamodb item: %w", err)
	}
	if len(output.Item) == 0 {
		return false, nil
	}
	if err := attributevalue.UnmarshalMap(output.Item, out); err != nil {
		return false, fmt.Errorf("unmarshal dynamodb item: %w", err)
	}
	return true, nil
}

func normalizeEnvPrefix(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return ""
	}
	return strings.TrimSuffix(prefix, "_") + "_"
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
