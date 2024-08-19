package domain

import "time"

// SignIn struct represents the data structure for signing in.
type SignIn struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:50;not null;unique" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"password_hash"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	UserID    uint       `gorm:"not null" json:"user_id"`  // Foreign key referencing the users table
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// SignInRepository interface defines the methods that any
// data storage provider needs to implement to get, create, and update SignIn.
type SignInRepository interface {
	CreateSignIn(signIn *SignIn) error
	GetSignInByUsername(username string) (*SignIn, error)
	GetSignInByEmail(email string) (*SignIn, error)
	UpdateSignIn(signIn *SignIn) error
	DeleteSignIn(id uint) error  
}