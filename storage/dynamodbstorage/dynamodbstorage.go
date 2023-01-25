// package dynamodbstorage provides a client for storing shortlinks in Amazon
// DynamoDB.
//
// The primary key is `pk` and the sort key is `sk`.
//
// Shortlinks have a `pk` of "s", and `sk` of their From value.
//
// Deleted shortlinks are exactly the same but with a `pk` of "d".
//
// History (previous versions of shortlinks) have a `pk` of "h" with their From
// value appended (ie the history of the "frew" shortlink has a `pk` of
// "hfrew") and an `sk` of the RFC3339 representation of the time that history
// was created.
package dynamodbstorage

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/frioux/shortlinks/shortlinks"
)

const (
	pkShortlink = "s"
	pkDeletedShortlink = "d"
)

type Client struct {
	DB *dynamodb.Client

	Table string
}

type shortlink struct {
	// PK is hardcoded to s for shortlinks, or d for deleted shortlinks.
	PK   string `dynamodbav:"pk"`
	From string `dynamodbav:"sk"`
	To   string `dynamodbav:"to,omitempty"`

	Description string `dynamodbav:"d,omitempty"`
}

func mustMarshal(v interface{}) map[string]types.AttributeValue {
	av, err := attributevalue.MarshalMap(v)
	if err != nil {
		panic(err)
	}

	return av
}

func mustUnmarshal(av map[string]types.AttributeValue, d interface{}) {
	if err := attributevalue.UnmarshalMap(av, d); err != nil {
		panic(err)
	}
}

func (cl *Client) Shortlink(from string) (shortlinks.Shortlink, error) {
	gio, err := cl.DB.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String(cl.Table),
		Key:       mustMarshal(shortlink{PK: pkShortlink, From: from}),
	})
	if err != nil {
		return shortlinks.Shortlink{}, err
	}

	var s shortlink
	mustUnmarshal(gio.Item, &s)

	return shortlinks.Shortlink{
		From: s.From,
		To:   s.To,

		Description: s.Description,
	}, nil
}

func (cl *Client) CreateShortlink(sl shortlinks.Shortlink) error {
	if _, err := cl.DB.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(cl.Table),
		Item: mustMarshal(shortlink{
			PK:   pkShortlink,
			From: sl.From,
			To:   sl.To,

			Description: sl.Description,
		}),
	}); err != nil {
		return err
	}

	return nil
}

func (cl *Client) pkShortlinks(pk string) ([]shortlinks.Shortlink, error) {
	qi := &dynamodb.QueryInput{
		TableName:              aws.String(cl.Table),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: pk},
		},
	}
	pager := dynamodb.NewQueryPaginator(cl.DB, qi)

	ret := make([]shortlinks.Shortlink, 0, 100)
	for pager.HasMorePages() {
		o, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, err
		}

		for _, itm := range o.Items {
			var s shortlink
			mustUnmarshal(itm, &s)
			ret = append(ret, shortlinks.Shortlink{
				From: s.From,
				To:   s.To,

				Description: s.Description,
			})
		}
	}

	return ret, nil
}

func (cl *Client) AllShortlinks() ([]shortlinks.Shortlink, error) { return cl.pkShortlinks(pkShortlink) }

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

	if _, err := cl.DB.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(cl.Table),
		Item: mustMarshal(shortlink{
			PK:   pkDeletedShortlink,
			From: sl.From,
			To:   sl.To,

			Description: sl.Description,
		}),
	}); err != nil {
		return err
	}

	_, err = cl.DB.DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
		TableName: aws.String(cl.Table),
		Key:       mustMarshal(shortlink{PK: pkShortlink, From: from}),
	})
	if err != nil {
		return err
	}

	return nil
}

func (cl *Client) DeletedShortlinks() ([]shortlinks.Shortlink, error) { return cl.pkShortlinks(pkDeletedShortlink) }

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
		TableName:              aws.String(cl.Table),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "h" + from},
		},
	}
	pager := dynamodb.NewQueryPaginator(cl.DB, qi)

	ret := make([]shortlinks.History, 0, 100)
	for pager.HasMorePages() {
		o, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, err
		}

		for _, itm := range o.Items {
			var h history
			mustUnmarshal(itm, &h)
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
	if _, err := cl.DB.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(cl.Table),
		Item: mustMarshal(history{
			PK:   "h" + h.From,
			When: time.Now().String(),
			Who:  h.Who,
			To:   h.To,

			Description: h.Description,
		}),
	}); err != nil {
		return err
	}

	return nil
}
