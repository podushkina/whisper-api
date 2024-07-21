// internal/model/auth.go
package model

type User struct {
	ID     string `json:"id"`
	APIKey string `json:"-"`
}
