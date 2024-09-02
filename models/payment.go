package models

import (
	"errors"
	"sync"
	"time"
)

type PaymentMode string

const (
	Cash         PaymentMode = "Cash"
	BankTransfer PaymentMode = "BankTransfer"
	UPI          PaymentMode = "UPI"
)

var (
	paymentIDCounter int32
	muLock           sync.Mutex // to ensure thread safety if accessed by multiple goroutines
)

// Payment struct represents a payment made to settle expenses.
type Payment struct {
	Group              string      // Name or ID of the group associated with the payment
	ID                 int32       // Unique identifier for the payment
	Payer              *User       // The user making the payment
	Payee              *User       // The user receiving the payment
	Amount             float64     // Amount paid
	Mode               PaymentMode // Mode of payment
	PaymentDate        time.Time   // Date and time of payment
	PaymentIdentifier  string      // Unique identifier for the payment (e.g., transaction ID)
	Note               string      // Additional notes for the payment
	AssociatedExpenses []*Expense  // List of expenses being settled by this payment
}

// generatePaymentID generates a unique ID for the payment.
func generatePaymentID() int32 {
	mu.Lock()
	defer mu.Unlock()
	paymentIDCounter++
	return paymentIDCounter
}

// NewPayment creates a new Payment instance.
func NewPayment(payer *User, payee *User, amount float64, mode PaymentMode, identifier string, note string, expenses []*Expense) *Payment {
	return &Payment{
		ID:                 generatePaymentID(),
		Payer:              payer,
		Payee:              payee,
		Amount:             amount,
		Mode:               mode,
		PaymentDate:        time.Now(),
		PaymentIdentifier:  identifier,
		Note:               note,
		AssociatedExpenses: expenses,
	}
}

func (p *Payment) SettlePayment() error {
	if p.Amount <= 0 {
		return errors.New("payment amount must be greater than zero")
	}

	totalExpenseAmount := 0.0
	for _, expense := range p.AssociatedExpenses {
		totalExpenseAmount += expense.Amount
	}

	if p.Amount < totalExpenseAmount {
		return errors.New("partial payment made, some expenses are still unsettled")
	}

	if p.Amount > totalExpenseAmount {
		return errors.New("payment amount exceeds the total amount of expenses")
	}

	// Assume logic to mark expenses as settled here
	//for _, expense := range p.AssociatedExpenses {
	//	// Mark each expense as settled (this is just a placeholder)
	//	// You may need to update the expense status or similar logic
	//}

	return nil
}
