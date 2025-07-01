package models

type PrivacyPolicy struct {
	Id       string `json:"id"`
	Content  string `json:"content"`
	Created  uint   `json:"created" db:"created"`
	Modified uint   `json:"modified" db:"modified"`
}
