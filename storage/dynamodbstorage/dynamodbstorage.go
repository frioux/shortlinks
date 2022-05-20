package dynamodbstorage

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/frioux/shortlinks/shortlinks"
)

type Client struct {
	db *dynamodb.Client

	table string
}

func NewClient() (*Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		return nil, err
	}

	svc := dynamodb.NewFromConfig(cfg)

	return &Client{
		db:    svc,
		table: "dev-zrorg--shortlinks",
	}, nil
}

type shortlink struct {
	// PK is hardcoded to s for shortlinks, or d for deleted shortlinks.
	PK   string `dynamodbav:"pk"`
	From string `dynamodbav:"sk"`
	To   string `dynamodbav:"to,omitempty"`

	Description string `dynamodbav:"d,omitempty"`
}

func (cl *Client) Shortlink(from string) (shortlinks.Shortlink, error) {
	var ret shortlinks.Shortlink

	av, err := attributevalue.MarshalMap(shortlink{PK: "s", From: from})
	if err != nil {
		return ret, err
	}

	gio, err := cl.db.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String(cl.table),
		Key:       av,
	})
	if err != nil {
		return ret, err
	}

	var s shortlink
	if err := attributevalue.UnmarshalMap(gio.Item, &s); err != nil {
		return ret, err
	}

	return shortlinks.Shortlink{
		From: s.From,
		To:   s.To,

		Description: s.Description,
	}, nil
}

func (cl *Client) CreateShortlink(sl shortlinks.Shortlink) error {
	av, err := attributevalue.MarshalMap(shortlink{
		PK:   "s",
		From: sl.From,
		To:   sl.To,

		Description: sl.Description,
	})
	if err != nil {
		return err
	}

	if _, err := cl.db.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(cl.table),
		Item:      av,
	}); err != nil {
		return err
	}

	return nil
}

func (cl *Client) pkShortlinks(pk string) ([]shortlinks.Shortlink, error) {
	qi := &dynamodb.QueryInput{
		TableName:              aws.String(cl.table),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: pk},
		},
	}
	pager := dynamodb.NewQueryPaginator(cl.db, qi)

	ret := make([]shortlinks.Shortlink, 0, 100)
	for pager.HasMorePages() {
		o, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, err
		}

		for _, itm := range o.Items {
			var s shortlink
			if err := attributevalue.UnmarshalMap(itm, &s); err != nil {
				return ret, err
			}
			ret = append(ret, shortlinks.Shortlink{
				From: s.From,
				To:   s.To,

				Description: s.Description,
			})
		}
	}

	return ret, nil
}

func (cl *Client) AllShortlinks() ([]shortlinks.Shortlink, error) { return cl.pkShortlinks("s") }

func (cl *Client) DeleteShortlink(from, who string) error {
	sl, err := cl.Shortlink(from)
	if err != nil {
		return err
	}

	if err := cl.InsertHistory(shortlinks.History{
		From: from,
		To:   sl.To,
		Who:  who,

		Description: "«deleted»",
	}); err != nil {
		return err
	}

	av, err := attributevalue.MarshalMap(shortlink{
		PK:   "d",
		From: sl.From,
		To:   sl.To,

		Description: sl.Description,
	})
	if err != nil {
		return err
	}

	if _, err := cl.db.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(cl.table),
		Item:      av,
	}); err != nil {
		return err
	}

	av, err = attributevalue.MarshalMap(shortlink{PK: "s", From: from})
	if err != nil {
		return err
	}
	_, err = cl.db.DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
		TableName: aws.String(cl.table),
		Key:       av,
	})
	if err != nil {
		return err
	}

	return nil
}

func (cl *Client) DeletedShortlinks() ([]shortlinks.Shortlink, error) { return cl.pkShortlinks("d") }

type history struct {
	// PK is h (for history) followed by the From value
	PK string `dynamodbav:"pk"`

	// When is the sk (RFC3339)
	When string `dynamodbav:"sk"`
	Who  string `dynamodbav:"who"`
	To   string `dynamodbav:"to,omitempty"`

	Description string `dynamodbav:"d,omitempty"`
}

func (h history) From() string {
	return strings.TrimPrefix(h.PK, "h")
}

func (cl *Client) History(from string) ([]shortlinks.History, error) {
	qi := &dynamodb.QueryInput{
		TableName:              aws.String(cl.table),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "h" + from},
		},
	}
	pager := dynamodb.NewQueryPaginator(cl.db, qi)

	ret := make([]shortlinks.History, 0, 100)
	for pager.HasMorePages() {
		o, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, err
		}

		for _, itm := range o.Items {
			var h history
			if err := attributevalue.UnmarshalMap(itm, &h); err != nil {
				return ret, err
			}
			ret = append(ret, shortlinks.History{
				From: h.From(),
				To:   h.To,
				When: h.When,
				Who:  h.Who,

				Description: h.Description,
			})
		}
	}

	return ret, nil
}

func (cl *Client) InsertHistory(h shortlinks.History) error {
	av, err := attributevalue.MarshalMap(history{
		PK:   "h" + h.From,
		When: time.Now().String(),
		Who:  h.Who,
		To:   h.To,

		Description: h.Description,
	})
	if err != nil {
		return err
	}

	if _, err := cl.db.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(cl.table),
		Item:      av,
	}); err != nil {
		return err
	}

	return nil
}
