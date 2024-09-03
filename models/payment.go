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
	Expenses   []*Expense // This is the correct field
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
		Timestamp:  time.Now(),
		Identifier: identifier,
		Note:       note,
		Expenses:   expenses,
	}
}

func (p *Payment) SettlePayment() error {
	if p.Amount <= 0 {
		return errors.New("payment amount must be greater than zero")
	}

	// Calculate the total amount of all expenses
	totalExpenseAmount := 0.0
	for _, expense := range p.Expenses {
		totalExpenseAmount += expense.Amount
	}
	fmt.Println("Total Expenses:", totalExpenseAmount)

	// Check if the payment amount is less than the total expense amount
	if p.Amount < totalExpenseAmount {
		return errors.New("payment amount is less than total expense amount")
	}
	fmt.Println("p.Amount:", p.Amount)

	remainingAmount := p.Amount

	for _, expense := range p.Expenses {
		fmt.Println("Expense is: ", expense)
		if expense.RemainingAmount <= 0 {
			fmt.Println("Remaining Expense Amount Is: ", expense.RemainingAmount)
			continue // Skip already settled expenses
		}

		// Calculate the payer's share of the expense
		payerShare := 0.0
		for i, user := range expense.SplitBetween {
			if user.Id == p.Payer.Id {
				payerShare = float64(expense.SplitRate[i]) / float64(len(expense.SplitBetween)) * expense.Amount
				break
			}
		}

		if payerShare > 0 {
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

	fmt.Println("Remaining Amount is ", remainingAmount)

	if remainingAmount > 0 {
		return errors.New("partial payment made, some expenses are still unsettled")
	}

	return nil
}
