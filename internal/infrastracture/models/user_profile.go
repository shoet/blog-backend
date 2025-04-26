package models

import "time"

type UserProfileId int64

type UserProfile struct {
	UserProfileId       UserProfileId `json:"userProfileId" db:"id"`
	UserId              UserId        `json:"userId" db:"user_id"`
	Nickname            string        `json:"nickname" db:"nickname"`
	AvatarImageFileName *string       `json:"avatarImageFileName,omitempty" db:"avatar_image_file_name"`
	AvatarImageFileURL  *string       `json:"avatarImageFileURL"`
	Biography           *string       `json:"bio" db:"bio"`
	Created             time.Time     `json:"created" db:"created"`
	Modified            time.Time     `json:"modified" db:"modified"`
}
