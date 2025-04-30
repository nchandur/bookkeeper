package model

type Author struct {
	Name        string `json:"name" bson:"name"`
	AuthorStats Stats  `json:"stats" bson:"stats"`
}

type Stats struct {
	Stars   float64 `json:"stars" bson:"stars"`
	Ratings float64 `json:"ratings" bson:"ratings"`
	Reviews float64 `json:"reviews" bson:"reviews"`
}
