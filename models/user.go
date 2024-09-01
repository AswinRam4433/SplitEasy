package models

import "sync"

var (
	userIDCounter int32
	mu            sync.Mutex // to ensure thread safety if accessed by multiple goroutines
)

type User struct {
	Name    string
	Balance float64
	Id      int32
}

func NewUser(name string) *User {
	mu.Lock()
	userIDCounter++
	id := userIDCounter
	mu.Unlock()
	return &User{
		Name:    name,
		Balance: 0.0,
		Id:      id,
	}
}
