package models

import "github.com/jmoiron/sqlx"

// Models struct that holds references to individual model structs
type Models struct {
	Users  UserModel  // Add other models as needed
	Tokens TokenModel // Add the TokenModel
}

// NewModels initializes all models and returns a Models struct
func NewModels(db *sqlx.DB) Models {
	return Models{
		Users:  UserModel{DB: db},
		Tokens: TokenModel{DB: db},
	}
}
