package model

type Subscription struct {
	ID       int
	UserID   int
	Plan     string
	IsActive bool
}
