package models

import (
	"errors"
	"fmt"
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
	ID         int
	Payer      *User
	Payee      *User
	Amount     float64
	Mode       PaymentMode
	Timestamp  time.Time
	Identifier string
	Note       string
	Expenses   []*Expense
}

// generatePaymentID generates a unique ID for the payment.
func generatePaymentID() int32 {
	muLock.Lock()
	defer muLock.Unlock()
	paymentIDCounter++
	return paymentIDCounter
}

// NewPayment creates a new Payment instance.
func NewPayment(payer *User, payee *User, amount float64, mode PaymentMode, identifier string, note string, expenses []*Expense) *Payment {
	return &Payment{
		ID:         int(generatePaymentID()),
		Payer:      payer,
		Payee:      payee,
		Amount:     amount,
		Mode:       mode,
		Timestamp:  time.Now().Truncate(time.Second), // Truncate to seconds for consistent comparison in tests
		Identifier: identifier,
		Note:       note,
		Expenses:   expenses,
	}
}

// SettlePayment settles the expenses based on the payment made.
func (p *Payment) SettlePayment() error {
	if p.Amount <= 0 {
		return errors.New("payment amount must be greater than zero")
	}

	remainingAmount := p.Amount

	for _, expense := range p.Expenses {
		if expense.RemainingAmount <= 0 {
			continue // Skip already settled expenses
		}

		// Calculate the payer's share of the expense
		payerShare := 0.0
		for i, user := range expense.SplitBetween {
			if user.Id == p.Payer.Id {
				// Share is calculated proportionally
				payerShare = float64(expense.SplitRate[i]) * expense.Amount
				break
			}
		}

		// Settle the payer's share
		if payerShare > 0 {
			// adjusting for payer's share of the expense
			if remainingAmount >= payerShare {
				remainingAmount -= payerShare
				expense.RemainingAmount -= payerShare
				p.Payee.Balance += payerShare
				p.Payer.Balance -= payerShare
			} else {
				expense.RemainingAmount -= remainingAmount
				p.Payee.Balance += remainingAmount
				p.Payer.Balance -= remainingAmount
				remainingAmount = 0
				break
			}
			expense.Payments = append(expense.Payments, p)
		}
	}

	if remainingAmount > 0 {
		return errors.New("partial payment made, some expenses are still unsettled")
	}

	return nil
}

// PrintPaymentInfo returns a formatted string containing all the fields of a Payment structure.
func PrintPaymentInfo(payment *Payment) string {
	expenseInfo := ""
	for _, expense := range payment.Expenses {
		expenseInfo += fmt.Sprintf("Expense ID: %d, Amount: %.2f, Paid By: %s, Remaining Amount: %.2f\n",
			expense.ID, expense.Amount, expense.PaidBy.Name, expense.RemainingAmount)
	}

	return fmt.Sprintf(
		"Payment Info:\nID: %d\nPayer: %s\nPayee: %s\nAmount: %.2f\nMode: %s\nTimestamp: %s\nIdentifier: %s\nNote: %s\nExpenses:\n%s",
		payment.ID,
		payment.Payer.Name,
		payment.Payee.Name,
		payment.Amount,
		payment.Mode,
		payment.Timestamp.Format(time.RFC3339),
		payment.Identifier,
		payment.Note,
		expenseInfo,
	)
}
