package recommend

import (
	"bookkeeper/internal/db"
	"context"
	"fmt"
	"log"
	"time"
)

func RecommendBooks(title string, topK int) {

	db.Connect()
	defer db.Client.Disconnect(context.TODO())

	collection := db.Client.Database("booksV2").Collection("works")

	start := time.Now()

	_, topDocs, err := GetTopKDocuments(collection, title, topK)

	if err != nil {
		log.Fatal(err)
	}

	duration := time.Since(start)

	for _, doc := range topDocs {
		fmt.Printf("%s (%f)\n", doc.Doc.Work.Title, doc.Score)
	}
	fmt.Println("\n\nBooks Retrieved in ", duration)

}
