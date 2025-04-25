package model

import "fmt"

type Book struct {
	ID        int       `bson:"book_id"`
	Title     string    `bson:"title"`
	Author    string    `bson:"author"`
	Summary   string    `bson:"summary"`
	Genres    []string  `bson:"genres"`
	Stars     float64   `bson:"stars"`
	Ratings   int       `bson:"ratings"`
	Reviews   int       `bson:"reviews"`
	Format    Format    `bson:"format"`
	Published string    `bson:"published"`
	URL       string    `bson:"url"`
	Embedding []float64 `bson:"embedding"`
}

type Format struct {
	PageNo int    `bson:"page_no"`
	Type   string `bson:"type"`
}

func (b *Book) Display() {
	fmt.Printf("BookID: %d\nTitle: %s\nAuthor: %s\nSummary:\n%s\nGenres: %v\nStars: %f\nRatings: %d\nReviews: %d\nPageCount: %d\nType: %s\nPublished: %v\nURL: %s\nEmedding: %v\n", b.ID, b.Title, b.Author, b.Summary, b.Genres, b.Stars, b.Ratings, b.Reviews, b.Format.PageNo, b.Format.Type, b.Published, b.URL, b.Embedding[:10])
}
