package models

import (
	"reflect"
	"sort"
	"testing"
)

func TestExpense_SplitExpense(t *testing.T) {
	type fields struct {
		Amount       float64
		PaidBy       *User
		SplitBetween []*User
		SplitRate    []float32
	}

	tests := []struct {
		name   string
		fields fields
		want   []struct {
			Id      int32
			Name    string
			Balance float64
		} // Expected balances of users after splitting
	}{
		{
			name: "Split between A, B, C with A paying",
			fields: fields{
				Amount:       300,
				PaidBy:       &User{Id: 1, Name: "A", Balance: 0},                                                                     // User A
				SplitBetween: []*User{{Id: 1, Name: "A", Balance: 0}, {Id: 2, Name: "B", Balance: 0}, {Id: 3, Name: "C", Balance: 0}}, // Users A, B, C
				SplitRate:    []float32{1, 1, 1},
			},
			want: []struct {
				Id      int32
				Name    string
				Balance float64
			}{
				{Id: 1, Name: "A", Balance: 200}, // A pays 300, split equally between A, B, C
				{Id: 2, Name: "B", Balance: -100},
				{Id: 3, Name: "C", Balance: -100},
			},
		},
		{
			name: "Valid Equal Split",
			fields: fields{
				Amount:       100,
				PaidBy:       &User{Id: 1, Name: "A", Balance: 0},                                     // User A
				SplitBetween: []*User{{Id: 1, Name: "A", Balance: 0}, {Id: 2, Name: "B", Balance: 0}}, // Users A, B
				SplitRate:    []float32{1, 1},
			},
			want: []struct {
				Id      int32
				Name    string
				Balance float64
			}{
				{Id: 1, Name: "A", Balance: 50}, // Split equally
				{Id: 2, Name: "B", Balance: -50},
			},
		},
		{
			name: "Valid Unequal Split",
			fields: fields{
				Amount:       150,
				PaidBy:       &User{Id: 1, Name: "A", Balance: 0},                                     // User A
				SplitBetween: []*User{{Id: 2, Name: "B", Balance: 0}, {Id: 3, Name: "C", Balance: 0}}, // Users B, C
				SplitRate:    []float32{1, 2},
			},
			want: []struct {
				Id      int32
				Name    string
				Balance float64
			}{
				{Id: 2, Name: "B", Balance: -50},
				{Id: 3, Name: "C", Balance: -100},
				{Id: 1, Name: "A", Balance: 150}, // Split unequally
			},
		},
		{
			name: "Amount Zero",
			fields: fields{
				Amount:       0,
				PaidBy:       &User{Id: 1, Name: "A", Balance: 100},                                     // User A
				SplitBetween: []*User{{Id: 2, Name: "B", Balance: 50}, {Id: 3, Name: "C", Balance: 50}}, // Users B, C
				SplitRate:    []float32{1, 1},
			},
			want: []struct {
				Id      int32
				Name    string
				Balance float64
			}{
				{Id: 1, Name: "A", Balance: 100}, // No change since amount is 0
				{Id: 2, Name: "B", Balance: 50},
				{Id: 3, Name: "C", Balance: 50},
			},
		},
		{
			name: "Nil PaidBy User",
			fields: fields{
				Amount:       100,
				PaidBy:       nil,                                                                       // No payer
				SplitBetween: []*User{{Id: 2, Name: "B", Balance: 50}, {Id: 3, Name: "C", Balance: 50}}, // Users B, C
				SplitRate:    []float32{1, 1},
			},
			want: []struct {
				Id      int32
				Name    string
				Balance float64
			}{
				{Id: 2, Name: "B", Balance: 50}, // No change since PaidBy is nil
				{Id: 3, Name: "C", Balance: 50},
			},
		},
		{
			name: "Empty SplitBetween",
			fields: fields{
				Amount:       100,
				PaidBy:       &User{Id: 1, Name: "A", Balance: 100}, // User A
				SplitBetween: []*User{},                             // No users to split
				SplitRate:    []float32{},
			},
			want: []struct {
				Id      int32
				Name    string
				Balance float64
			}{
				{Id: 1, Name: "A", Balance: 100}, // No change since no users to split between
			},
		},
		{
			name: "Different Split Rates",
			fields: fields{
				Amount:       300,
				PaidBy:       &User{Id: 1, Name: "A", Balance: 0},                                     // User A
				SplitBetween: []*User{{Id: 2, Name: "B", Balance: 0}, {Id: 3, Name: "C", Balance: 0}}, // Users B, C
				SplitRate:    []float32{1, 2},                                                         // B's share is half of C's
			},
			want: []struct {
				Id      int32
				Name    string
				Balance float64
			}{
				{Id: 2, Name: "B", Balance: -100},
				{Id: 3, Name: "C", Balance: -200},
				{Id: 1, Name: "A", Balance: 300}, // B pays 100, C pays 200, A receives 300
				// A should only be present once with the final balance
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Expense{
				Amount:       tt.fields.Amount,
				PaidBy:       tt.fields.PaidBy,
				SplitBetween: tt.fields.SplitBetween,
				SplitRate:    tt.fields.SplitRate,
			}

			e.SplitExpense()

			got := []struct {
				Id      int32
				Name    string
				Balance float64
			}{}

			userMap := make(map[int32]struct {
				Id      int32
				Name    string
				Balance float64
			})

			if e.PaidBy != nil {
				userMap[e.PaidBy.Id] = struct {
					Id      int32
					Name    string
					Balance float64
				}{Id: e.PaidBy.Id, Name: e.PaidBy.Name, Balance: e.PaidBy.Balance}
			}
			for _, user := range e.SplitBetween {
				userMap[user.Id] = struct {
					Id      int32
					Name    string
					Balance float64
				}{Id: user.Id, Name: user.Name, Balance: user.Balance}
			}

			for _, user := range userMap {
				got = append(got, user)
			}

			// Sort the `got` and `want` slices by ID to ensure order does not affect comparison
			sort.Slice(got, func(i, j int) bool {
				return got[i].Id < got[j].Id
			})
			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i].Id < tt.want[j].Id
			})

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expense.SplitExpense() = %v, want %v", got, tt.want)
			}
		})
	}
}
