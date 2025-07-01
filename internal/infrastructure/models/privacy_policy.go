package models

import "time"

type PrivacyPolicy struct {
	Id       string    `json:"id"`
	Content  string    `json:"content"`
	Created  time.Time `json:"created" db:"created"`
	Modified time.Time `json:"modified" db:"modified"`
}
