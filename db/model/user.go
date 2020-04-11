package model

type User struct {
	UserID       string `db:"user_id"`
	UserName     string `db:"username"`
	Name         string `db:"name"`
	FirstName    string `db:"firstname"`
	LastName     string `db:"lastname"`
	Email        string `db:"email"`
	FullPhotoURL string `db:"full_photo_url"`
	AccessToken  string `db:"access_token"`
	RefreshToken string `db:"refresh_token"`
	InstanceURL  string `db:"instance_url"`
	SfIsActive   bool   `db:"sf_is_active"`
	SfIsFrozen   bool   `db:"sf_is_frozen"`
}

type Users []User
