package models

type Movie struct {
	Title  string `json:"Title,omitempty" bson:"title,omitempty"`
	Year   string `json:"Year,omitempty"  bson:"year,omitempty"`
	Poster string `json:"Poster,omitempty" bson:"poster,omitempty"`
	ImdbID string `json:"imdbID,omitempty" bson:"imdb_id,omitempty"`
}
