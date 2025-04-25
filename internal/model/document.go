package model

type Document struct {
	Source string `bson:"source"`
	Work   Book   `bson:"work"`
}

type ScoredDocument struct {
	Doc   Document `bson:"document"`
	Score float64  `bson:"score"`
}
