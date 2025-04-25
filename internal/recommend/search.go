package recommend

import (
	"bookkeeper/internal/model"
	"container/heap"
	"context"
	"fmt"
	"math"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// get document by matching title
func GetDocument(collection *mongo.Collection, title string) (model.Document, error) {
	var res model.Document

	filter := bson.M{"work.title": bson.M{
		"$regex":   title,
		"$options": "i",
	}}

	opts := options.Find()
	opts.SetSort(bson.D{
		{Key: "work.ratings", Value: -1},
		{Key: "work.stars", Value: -1},
	})

	err := collection.FindOne(context.Background(), filter).Decode(&res)

	if err != nil {
		return res, err
	}

	return res, nil
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

	input, err := GetDocument(collection, title)

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
