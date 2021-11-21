package user

import "time"

type User struct {
	ID          string    `gorm:"column:user_id; not null, primaryKey"`
	Email       string    `gorm:"column:email; not null"`
	Username    string    `gorm:"column:username; not null"`
	DisplayName string    `gorm:"column:display_name; not null"`
	Biography   string    `gorm:"column:biography;"`
	ExternalUrl string    `gorm:"column:external_url"`
	IsVerified  bool      `gorm:"column:is_verified;"`
	Password    string    `gorm:"column:password; not null"`
	CreatedAt   time.Time `gorm:"column:created_at; not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at; not null"`
}

func (u *User) ToResponse() *Response {
	return &Response{
		Username:          u.Username,
		DisplayName:       u.DisplayName,
		Biography:         u.Biography,
		ExternalUrl:       u.ExternalUrl,
		ProfilePictureURL: u.Email,
	}
}
