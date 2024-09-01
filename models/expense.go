package models

type Expense struct {
	Amount       float64
	PaidBy       *User
	SplitBetween []*User
}

func NewExpense(amount float64, paidBy *User, splitBetween []*User) *Expense {
	return &Expense{
		Amount:       amount,
		PaidBy:       paidBy,
		SplitBetween: splitBetween,
	}
}

func (e *Expense) SplitExpense() {
	splitAmount := e.Amount / float64(len(e.SplitBetween))
	for _, user := range e.SplitBetween {
		if user == e.PaidBy {
			user.Balance += e.Amount - splitAmount
		} else {
			user.Balance -= splitAmount
		}
	}
}
