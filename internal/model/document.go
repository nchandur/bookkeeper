package model

type Document struct {
	Source string `json:"source" bson:"source"`
	Work   Book   `json:"work" bson:"work"`
}

type ScoredDocument struct {
	Doc   Document `json:"document" bson:"document"`
	Score float64  `json:"score" bson:"score"`
}
