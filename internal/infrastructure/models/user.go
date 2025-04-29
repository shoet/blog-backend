package models

type UserId int64

type User struct {
	Id       UserId       `json:"id,omitempty" db:"id"`
	Name     string       `json:"name" db:"name"`
	Email    string       `json:"email,omitempty" db:"email"`
	Password string       `json:"password,omitempty" db:"password"`
	Profile  *UserProfile `json:"profile,omitempty"`
	Created  uint         `json:"created,omitempty" db:"created"`
	Modified uint         `json:"modified,omitempty" db:"modified"`
}
