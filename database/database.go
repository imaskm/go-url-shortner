package database

import (
	"context"
	"encoding/json"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	conn       *mongo.Client
	databse    string
	collection string
}

type URL struct {
	LongUrl  string `bson:"long_url" json:"long_url"`
	ShortUrl string `bson:"short_url" json:"short_url"`
}

func NewDatabase() *MongoDB {

	opts := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Fatal(err)
	}

	return &MongoDB{
		conn:       client,
		databse:    "url-shortner",
		collection: "urls",
	}

}

func (m *MongoDB) SaveShortURL(longURL, shortURL string) error {

	_, err := m.conn.Database(m.databse).Collection(m.collection).InsertOne(context.Background(),
		bson.M{
			"short_url": shortURL,
			"long_url":  longURL,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *MongoDB) GetShortURLForLongURL(longURL string) (string, error) {
	result := bson.M{}

	opts := options.FindOne().SetProjection(bson.M{
		"short_url": 1,
		"_id":       0,
		"long_url":  1,
	})

	err := m.conn.Database(m.databse).Collection(m.collection).FindOne(
		context.Background(),
		bson.M{
			"long_url": longURL,
		}, opts,
	).Decode(&result)

	if err != nil {
		return "", nil
	}

	res := &URL{}

	raw, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(raw, res)
	if err != nil {
		return "", err
	}

	return res.ShortUrl, nil
}

func (m *MongoDB) GetLongURLForShortURL(shortURL string) (string, error) {
	result := bson.M{}

	opts := options.FindOne().SetProjection(bson.M{
		"short_url": 1,
		"_id":       0,
		"long_url":  1,
	})

	err := m.conn.Database(m.databse).Collection(m.collection).FindOne(
		context.Background(),
		bson.M{
			"short_url": shortURL,
		}, opts,
	).Decode(&result)

	if err != nil {
		return "", nil
	}

	res := &URL{}

	raw, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(raw, res)
	if err != nil {
		return "", err
	}

	return res.LongUrl, nil
}
