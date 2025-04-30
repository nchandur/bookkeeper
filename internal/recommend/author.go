package recommend

import (
	"bookkeeper/internal/model"
	"context"
	"fmt"
	"math"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getAuthorName(ctx context.Context, collection *mongo.Collection, author string) (string, error) {
	exactFilter := bson.M{
		"work.author": bson.M{
			"$regex":   fmt.Sprintf("^%s$", regexpEscape(author)),
			"$options": "i",
		},
	}
	projection := bson.M{"work.author": 1, "_id": 0}

	var result struct {
		Work struct {
			Author string `bson:"author"`
		} `bson:"work"`
	}

	err := collection.FindOne(ctx, exactFilter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err == nil {
		return result.Work.Author, nil
	}
	if err != mongo.ErrNoDocuments {
		return "", err
	}

	partialFilter := bson.M{
		"work.author": bson.M{
			"$regex":   regexpEscape(author),
			"$options": "i",
		},
	}

	err = collection.FindOne(ctx, partialFilter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("author not found")
		}
		return "", err
	}

	return result.Work.Author, nil
}

func regexpEscape(input string) string {
	specialChars := `.+*?()|[]{}^$`
	for _, c := range specialChars {
		input = strings.ReplaceAll(input, string(c), `\`+string(c))
	}
	return input
}

func GetAuthor(ctx context.Context, collection *mongo.Collection, author string) (model.Author, error) {
	author, err := getAuthorName(ctx, collection, author)

	if err != nil {
		return model.Author{}, err
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"work.author": author,
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":     nil,
			"stars":   bson.M{"$avg": "$work.stars"},
			"ratings": bson.M{"$avg": "$work.ratings"},
			"reviews": bson.M{"$avg": "$work.reviews"},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return model.Author{Name: author}, err
	}
	defer cursor.Close(ctx)

	var result model.Stats

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return model.Author{Name: author}, err
		}
	}

	result.Ratings = math.Ceil(result.Ratings)
	result.Reviews = math.Ceil(result.Reviews)

	result.Stars = math.Ceil(result.Stars*100) / 100

	return model.Author{Name: author, AuthorStats: result}, nil
}
