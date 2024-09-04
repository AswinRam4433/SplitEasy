package models

import (
	"testing"
	"time"
)

func TestNewPayment(t *testing.T) {
	type args struct {
		payer      *User
		payee      *User
		amount     float64
		mode       PaymentMode
		identifier string
		note       string
		expenses   []*Expense
	}
	tests := []struct {
		name string
		args args
		want *Payment
	}{
		{
			name: "Simple Payment",
			args: args{
				payer:      &User{Name: "User A", Balance: 200, Id: 1},
				payee:      &User{Name: "User B", Balance: 100, Id: 2},
				amount:     100, // Correct amount to settle both expenses
				mode:       UPI,
				identifier: "DUMMYTXN1",
				note:       "Lorem Ipsum",
				expenses:   []*Expense{{Amount: 100, PaidBy: &User{Name: "User A", Balance: 200, Id: 1}, SplitBetween: []*User{{Name: "User A", Balance: 200, Id: 1}, {Name: "User B", Balance: 100, Id: 2}}, SplitRate: []float32{1.0}}},
			},
			want: &Payment{
				Payer:      &User{Name: "User A", Balance: 200, Id: 1},
				Payee:      &User{Name: "User B", Balance: 100, Id: 2},
				Amount:     100,
				Mode:       UPI,
				Timestamp:  time.Now(), // Timestamp field
				Identifier: "DUMMYTXN1",
				Note:       "Lorem Ipsum",
				Expenses:   []*Expense{{Amount: 100, PaidBy: &User{Name: "User A", Balance: 200, Id: 1}, SplitBetween: []*User{{Name: "User A", Balance: 200, Id: 1}, {Name: "User B", Balance: 100, Id: 2}}, SplitRate: []float32{1.0}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the payment
			got := NewPayment(tt.args.payer, tt.args.payee, tt.args.amount, tt.args.mode, tt.args.identifier, tt.args.note, tt.args.expenses)

			// Check Timestamp separately
			if !got.Timestamp.Truncate(24 * time.Hour).Equal(tt.want.Timestamp.Truncate(24 * time.Hour)) {
				t.Errorf("NewPayment() Timestamp = %v, want %v", got.Timestamp, tt.want.Timestamp)
			}

			// Directly compare each field
			if got.Payer.Name != tt.want.Payer.Name || got.Payer.Balance != tt.want.Payer.Balance || got.Payer.Id != tt.want.Payer.Id {
				t.Errorf("NewPayment() Payer = %v, want %v", got.Payer, tt.want.Payer)
			}
			if got.Payee.Name != tt.want.Payee.Name || got.Payee.Balance != tt.want.Payee.Balance || got.Payee.Id != tt.want.Payee.Id {
				t.Errorf("NewPayment() Payee = %v, want %v", got.Payee, tt.want.Payee)
			}
			if got.Amount != tt.want.Amount {
				t.Errorf("NewPayment() Amount = %v, want %v", got.Amount, tt.want.Amount)
			}
			if got.Mode != tt.want.Mode {
				t.Errorf("NewPayment() Mode = %v, want %v", got.Mode, tt.want.Mode)
			}
			if got.Identifier != tt.want.Identifier {
				t.Errorf("NewPayment() Identifier = %v, want %v", got.Identifier, tt.want.Identifier)
			}
			if got.Note != tt.want.Note {
				t.Errorf("NewPayment() Note = %v, want %v", got.Note, tt.want.Note)
			}
		})
	}
}

func TestPayment_SettlePayment(t *testing.T) {
	payer := &User{Name: "User A", Balance: 0, Id: 1}
	payee := &User{Name: "User B", Balance: 0, Id: 2}
	expense1 := &Expense{
		Amount:          100,
		PaidBy:          payer,
		SplitBetween:    []*User{payer, payee},
		SplitRate:       []float32{1.0, 1.0},
		RemainingAmount: 100, // Amount remaining to be settled
	}
	expense2 := &Expense{
		Amount:          100,
		PaidBy:          payer,
		SplitBetween:    []*User{payer, payee},
		SplitRate:       []float32{1.0, 1.0},
		RemainingAmount: 100, // Amount remaining to be settled
	}
	expense3 := &Expense{
		Amount:          200,
		PaidBy:          payee,
		SplitBetween:    []*User{payer, payee},
		SplitRate:       []float32{1.0, 1.0},
		RemainingAmount: 200, // Amount remaining to be settled
	}

	type fields struct {
		Payer      *User
		Payee      *User
		Amount     float64
		Mode       PaymentMode
		Timestamp  time.Time
		Identifier string
		Note       string
		Expenses   []*Expense
	}

	tests := []struct {
		name   string
		fields fields
		// Expected state of the Payment object
		wantErr bool
	}{
		{
			name: "Settle Payment Success",
			fields: fields{
				Payer:    payer,
				Payee:    payee,
				Amount:   200,
				Mode:     Cash,
				Expenses: []*Expense{expense1, expense2},
			},
			wantErr: false,
		},
		{
			name: "Settle Payment Success #2",
			fields: fields{
				Payer:    payer,
				Payee:    payee,
				Amount:   150,
				Mode:     Cash,
				Expenses: []*Expense{expense1, expense3},
			},
			wantErr: false,
		},
		{
			name: "Excess Payment",
			fields: fields{
				Payer:    payer,
				Payee:    payee,
				Amount:   10000,
				Mode:     Cash,
				Expenses: []*Expense{expense1, expense3},
			},
			wantErr: true,
		},
		{
			name: "Negative Amount Payment ",
			fields: fields{
				Payer:    payer,
				Payee:    payee,
				Amount:   -10000,
				Mode:     Cash,
				Expenses: []*Expense{expense1, expense2},
			},
			wantErr: true,
		},
		{
			name: "Partial Payment Error",
			fields: fields{
				Payer:    payer,
				Payee:    payee,
				Amount:   50,
				Mode:     Cash,
				Expenses: []*Expense{expense1, expense2},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payment := &Payment{
				Payer:      tt.fields.Payer,
				Payee:      tt.fields.Payee,
				Amount:     tt.fields.Amount,
				Mode:       tt.fields.Mode,
				Identifier: tt.fields.Identifier,
				Note:       tt.fields.Note,
				Expenses:   tt.fields.Expenses,
			}
			if err := payment.SettlePayment(); (err != nil) != tt.wantErr {
				t.Errorf("SettlePayment() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
