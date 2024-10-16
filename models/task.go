package models

type Task struct {
	ID      string `json:"id,omitempty" bson:"_id,omitempty"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  string `json:"user_id"`
}
