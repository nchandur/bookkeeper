package model

import (
	"fmt"
	"time"
)

type Book struct {
	ID        int       `json:"book_id" bson:"book_id"`
	Title     string    `json:"title" bson:"title"`
	Author    string    `json:"author" bson:"author"`
	Summary   string    `json:"summary" bson:"summary"`
	Genres    []string  `json:"genres" bson:"genres"`
	Stars     float64   `json:"stars" bson:"stars"`
	Ratings   int       `json:"ratings" bson:"ratings"`
	Reviews   int       `json:"reviews" bson:"reviews"`
	Format    Format    `json:"format" bson:"format"`
	Published time.Time `json:"published" bson:"published"`
	URL       string    `json:"url" bson:"url"`
	Embedding []float64 `json:"embedding" bson:"embedding"`
}

type Format struct {
	PageNo int    `json:"page_no" bson:"page_no"`
	Type   string `json:"type" bson:"type"`
}

func (b *Book) Display() {
	fmt.Printf("BookID: %d\nTitle: %s\nAuthor: %s\nSummary:\n%s\nGenres: %v\nStars: %f\nRatings: %d\nReviews: %d\nPageCount: %d\nType: %s\nPublished: %v\nURL: %s\n", b.ID, b.Title, b.Author, b.Summary, b.Genres, b.Stars, b.Ratings, b.Reviews, b.Format.PageNo, b.Format.Type, b.Published, b.URL)
}
