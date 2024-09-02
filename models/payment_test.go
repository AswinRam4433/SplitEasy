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
				amount:     300,
				mode:       UPI,
				identifier: "DUMMYTXN1",
				note:       "Lorem Ipsum",
				expenses:   []*Expense{{Amount: 100, PaidBy: &User{Name: "User A", Balance: 200, Id: 1}, SplitBetween: []*User{{Name: "User B", Balance: 100, Id: 2}}, SplitRate: []float32{1.0}}},
			},
			want: &Payment{
				Payer:             &User{Name: "User A", Balance: 200, Id: 1},
				Payee:             &User{Name: "User B", Balance: 100, Id: 2},
				Amount:            300,
				Mode:              UPI,
				PaymentDate:       time.Date(2024, 9, 2, 14, 30, 45, 100, time.Local),
				PaymentIdentifier: "DUMMYTXN1",
				Note:              "Lorem Ipsum",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the payment
			got := NewPayment(tt.args.payer, tt.args.payee, tt.args.amount, tt.args.mode, tt.args.identifier, tt.args.note, tt.args.expenses)

			// Compare the ID separately if necessary
			if got.ID == 0 {
				// Check if the ID has been set correctly
				if got.ID <= 0 {
					t.Errorf("NewPayment() ID = %v, want > 0", got.ID)
				}
			}

			// Check PaymentDate separately
			if !got.PaymentDate.Truncate(24 * time.Hour).Equal(tt.want.PaymentDate.Truncate(24 * time.Hour)) {
				t.Errorf("NewPayment() PaymentDate = %v, want %v", got.PaymentDate, tt.want.PaymentDate)
			}

			if got.ID <= 0 {
				t.Errorf("NewPayment() ID = %v, want > 0", got.ID)
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
			if got.PaymentIdentifier != tt.want.PaymentIdentifier {
				t.Errorf("NewPayment() PaymentIdentifier = %v, want %v", got.PaymentIdentifier, tt.want.PaymentIdentifier)
			}
			if got.Note != tt.want.Note {
				t.Errorf("NewPayment() Note = %v, want %v", got.Note, tt.want.Note)
			}
		})
	}
}

func TestPayment_SettlePayment1(t *testing.T) {
	payer := &User{Name: "User A", Balance: -200, Id: 1}
	payee := &User{Name: "User B", Balance: 100, Id: 2}
	expense1 := &Expense{
		Amount:       100,
		PaidBy:       payee,
		SplitBetween: []*User{payer, payee},
		SplitRate:    []float32{1.0, 1.0},
	}
	expense2 := &Expense{
		Amount:       150,
		PaidBy:       payee,
		SplitBetween: []*User{payer, payee},
		SplitRate:    []float32{1.0, 1.0},
	}
	type fields struct {
		Group              string
		ID                 int32
		Payer              *User
		Payee              *User
		Amount             float64
		Mode               PaymentMode
		PaymentDate        time.Time
		PaymentIdentifier  string
		Note               string
		AssociatedExpenses []*Expense
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Settle Payment Success",
			fields: fields{
				Payer:              payer,
				Payee:              payee,
				Amount:             250,
				Mode:               Cash,
				AssociatedExpenses: []*Expense{expense1, expense2},
			},
			wantErr: false,
		},
		{
			name: "Partial Payment Error",
			fields: fields{
				Payer:              payer,
				Payee:              payee,
				Amount:             100, // Less than the total of expenses
				Mode:               BankTransfer,
				AssociatedExpenses: []*Expense{expense1, expense2},
			},
			wantErr: true,
		},
		{
			name: "Zero Amount Payment Error",
			fields: fields{
				Payer:              payer,
				Payee:              payee,
				Amount:             0, // Zero amount
				Mode:               UPI,
				AssociatedExpenses: []*Expense{expense1},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Payment{
				Group:              tt.fields.Group,
				ID:                 tt.fields.ID,
				Payer:              tt.fields.Payer,
				Payee:              tt.fields.Payee,
				Amount:             tt.fields.Amount,
				Mode:               tt.fields.Mode,
				PaymentDate:        tt.fields.PaymentDate,
				PaymentIdentifier:  tt.fields.PaymentIdentifier,
				Note:               tt.fields.Note,
				AssociatedExpenses: tt.fields.AssociatedExpenses,
			}
			if err := p.SettlePayment(); (err != nil) != tt.wantErr {
				t.Errorf("SettlePayment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPayment_SettlePayment(t *testing.T) {
	payer := &User{Name: "User A", Balance: -200, Id: 1}
	payee := &User{Name: "User B", Balance: 100, Id: 2}
	expense1 := &Expense{
		Amount:       100,
		PaidBy:       payee,
		SplitBetween: []*User{payer, payee},
		SplitRate:    []float32{1.0, 1.0},
	}
	expense2 := &Expense{
		Amount:       150,
		PaidBy:       payee,
		SplitBetween: []*User{payer, payee},
		SplitRate:    []float32{1.0, 1.0},
	}

	type fields struct {
		Group              string
		ID                 int32
		Payer              *User
		Payee              *User
		Amount             float64
		Mode               PaymentMode
		PaymentDate        time.Time
		PaymentIdentifier  string
		Note               string
		AssociatedExpenses []*Expense
	}

	tests := []struct {
		name    string
		fields  fields
		want    fields // Expected state of the Payment object
		wantErr bool
	}{
		{
			name: "Settle Payment Success",
			fields: fields{
				Payer:              payer,
				Payee:              payee,
				Amount:             250,
				Mode:               Cash,
				AssociatedExpenses: []*Expense{expense1, expense2},
			},
			want: fields{
				Payer:              payer,
				Payee:              payee,
				Amount:             250,
				Mode:               Cash,
				AssociatedExpenses: []*Expense{expense1, expense2},
			},
			wantErr: false,
		},
		{
			name: "Partial Payment Error",
			fields: fields{
				Payer:              payer,
				Payee:              payee,
				Amount:             100, // Less than the total of expenses
				Mode:               BankTransfer,
				AssociatedExpenses: []*Expense{expense1, expense2},
			},
			want: fields{
				Payer:              payer,
				Payee:              payee,
				Amount:             100, // Should be unchanged
				Mode:               BankTransfer,
				AssociatedExpenses: []*Expense{expense1, expense2},
			},
			wantErr: true,
		},
		{
			name: "Zero Amount Payment Error",
			fields: fields{
				Payer:              payer,
				Payee:              payee,
				Amount:             0, // Zero amount
				Mode:               UPI,
				AssociatedExpenses: []*Expense{expense1},
			},
			want: fields{
				Payer:              payer,
				Payee:              payee,
				Amount:             0, // Should be unchanged
				Mode:               UPI,
				AssociatedExpenses: []*Expense{expense1},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Payment{
				Group:              tt.fields.Group,
				ID:                 tt.fields.ID,
				Payer:              tt.fields.Payer,
				Payee:              tt.fields.Payee,
				Amount:             tt.fields.Amount,
				Mode:               tt.fields.Mode,
				PaymentDate:        tt.fields.PaymentDate,
				PaymentIdentifier:  tt.fields.PaymentIdentifier,
				Note:               tt.fields.Note,
				AssociatedExpenses: tt.fields.AssociatedExpenses,
			}
			err := p.SettlePayment()

			// Check if the error status is as expected
			if (err != nil) != tt.wantErr {
				t.Errorf("SettlePayment() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Compare Payment fields directly
			if p.Payer != tt.want.Payer || p.Payee != tt.want.Payee || p.Amount != tt.want.Amount || p.Mode != tt.want.Mode || p.PaymentIdentifier != tt.want.PaymentIdentifier || p.Note != tt.want.Note {
				t.Errorf("SettlePayment() = %+v, want %+v", p, tt.want)
			}
		})
	}
}
