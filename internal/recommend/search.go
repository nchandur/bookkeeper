package recommend

import (
	"bookkeeper/internal/model"
	"container/heap"
	"context"
	"fmt"
	"math"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// get document by matching title

func GetDocument(ctx context.Context, collection *mongo.Collection, title string) (model.Document, error) {
	exactRegex := fmt.Sprintf("^%s$", regexp.QuoteMeta(title))
	exactFilter := bson.M{
		"work.title": bson.M{
			"$regex":   exactRegex,
			"$options": "i",
		},
	}

	var exactResult model.Document
	err := collection.FindOne(ctx, exactFilter).Decode(&exactResult)
	if err == nil {
		return exactResult, nil
	}
	if err != mongo.ErrNoDocuments {
		return model.Document{}, err
	}

	substringFilter := bson.M{
		"work.title": bson.M{
			"$regex":   title,
			"$options": "i",
		},
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "work.title", Value: 1}}).
		SetLimit(1)

	cursor, err := collection.Find(ctx, substringFilter, opts)
	if err != nil {
		return model.Document{}, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var substringResult model.Document
		if err := cursor.Decode(&substringResult); err != nil {
			return model.Document{}, err
		}
		return substringResult, nil
	}

	return model.Document{}, mongo.ErrNoDocuments
}

// performs cosine similarity between two vectors
func cosineSimilarity(a, b []float64) (float64, error) {

	if len(a) == 0 || len(b) == 0 || len(a) != len(b) {
		return 0, fmt.Errorf("invalid vector")
	}
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0, nil
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB)), nil
}

// retrieve Top K documents that are closest to the given title
func GetTopKDocuments(collection *mongo.Collection, title string, topK int) (model.Book, []model.ScoredDocument, error) {
	ctx := context.TODO()

	pq := &PriorityQueue{}
	heap.Init(pq)

	if len(title) == 0 {
		return model.Book{}, nil, fmt.Errorf("empty title")
	}

	input, err := GetDocument(ctx, collection, title)

	if err != nil {
		return model.Book{}, nil, err
	}

	fmt.Println("retrieved book with title: ", input.Work.Title)

	var cur *mongo.Cursor

	cur, err = collection.Find(ctx, bson.M{})

	if err != nil {
		return model.Book{}, nil, err
	}

	count := 0

	for cur.Next(ctx) {
		var doc model.Document
		cur.Decode(&doc)

		if input.Work.Title != doc.Work.Title {
			cosSim, _ := cosineSimilarity(input.Work.Embedding, doc.Work.Embedding)

			InsertIntoQueue(pq, model.ScoredDocument{Doc: doc, Score: cosSim}, topK)

		}

		count++
		fmt.Printf("\r%d documents traversed...", count)
	}

	fmt.Printf("\n\n")
	cur.Close(ctx)

	topDocs := make([]model.ScoredDocument, pq.Len())

	for i := pq.Len() - 1; i >= 0; i-- {
		topDocs[i] = heap.Pop(pq).(model.ScoredDocument)
	}

	return input.Work, topDocs, nil

}
