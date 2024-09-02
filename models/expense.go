package models

import "errors"

type Expense struct {
	Amount       float64
	PaidBy       *User
	SplitBetween []*User
	SplitRate    []float32
}

func NewEqualExpense(amount float64, paidBy *User, splitBetween []*User) (*Expense, error) {
	splitRate := make([]float32, len(splitBetween))
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}
	if paidBy == nil {
		return nil, errors.New("paidBy cannot be nil")
	}
	if splitBetween == nil {
		return nil, errors.New("splitBetween cannot be nil")
	}
	if len(splitRate) != len(splitBetween) {
		return nil, errors.New("splitRate length must be equal to splitBetween")
	}
	for i := range splitRate {
		splitRate[i] = 1.0 // Default equal rate
	}

	return &Expense{
		Amount:       amount,
		PaidBy:       paidBy,
		SplitBetween: splitBetween,
		SplitRate:    splitRate,
	}, nil
}

func NewExpense(amount float64, paidBy *User, splitBetween []*User, splitRate []float32) (*Expense, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}
	if paidBy == nil {
		return nil, errors.New("paidBy cannot be nil")
	}
	if splitBetween == nil {
		return nil, errors.New("splitBetween cannot be nil")
	}
	if len(splitRate) != len(splitBetween) {
		return nil, errors.New("splitRate length must be equal to splitBetween")
	}
	return &Expense{
		Amount:       amount,
		PaidBy:       paidBy,
		SplitBetween: splitBetween,
		SplitRate:    splitRate,
	}, nil
}

func (e *Expense) SplitExpense() error {
	// Handle the case where PaidBy is nil; no operation should be performed
	if e.PaidBy == nil {
		return errors.New("paidBy cannot be nil")
	}

	if e.Amount == 0 {
		return nil // No change since amount is 0
	}

	totalSplitRate := 0.0
	// Calculate the total split rate
	for _, rate := range e.SplitRate {
		totalSplitRate += float64(rate)
	}

	// Check for division by zero
	if totalSplitRate == 0 {
		return nil // No rates set, nothing to split
	}

	// Create a map to keep track of balances for each user
	userBalances := make(map[int32]float64)
	for _, user := range e.SplitBetween {
		userBalances[user.Id] = 0
	}
	if e.PaidBy != nil {
		userBalances[e.PaidBy.Id] = e.PaidBy.Balance
	}

	// Deduct the shares from all users including the PaidBy user
	for i, user := range e.SplitBetween {
		splitAmount := (float64(e.SplitRate[i]) / totalSplitRate) * e.Amount
		userBalances[user.Id] -= splitAmount
	}

	// Update the balance of the PaidBy user by adding the total amount
	if e.PaidBy != nil {
		userBalances[e.PaidBy.Id] += e.Amount
	}

	// Update the balances in the users
	for _, user := range e.SplitBetween {
		user.Balance = userBalances[user.Id]
	}
	if e.PaidBy != nil {
		e.PaidBy.Balance = userBalances[e.PaidBy.Id]
	}

	return nil
}
