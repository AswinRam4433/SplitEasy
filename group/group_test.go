package group

import (
	"fmt"
	"reflect"
	"splitwise/models"
	"testing"
)

func TestGroup_AddMember(t *testing.T) {
	type fields struct {
		Name    string
		Members []*models.User
	}
	type args struct {
		user *models.User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []models.User
	}{
		{name: "Empty Group", fields: fields{Name: "Empty Group", Members: []*models.User{}}, args: args{models.NewUser("User A")}, want: []models.User{{Name: "User A", Balance: 0, Id: 1}}},
		{name: "Non Empty Group", fields: fields{Name: "Non Empty Group", Members: []*models.User{{Name: "User A", Balance: 0, Id: 1}}}, args: args{models.NewUser("User B")}, want: []models.User{{Name: "User A", Balance: 0, Id: 1}, {Name: "User B", Balance: 0, Id: 2}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				Name:    tt.fields.Name,
				Members: tt.fields.Members,
			}
			g.AddMember(tt.args.user)
			fmt.Println(g.ListMembers())
			if !reflect.DeepEqual(g.ListMembers(), tt.want) {
				t.Errorf("AddMember() = %v, want %v", g, tt.want)
			}
		})
	}
}

func TestGroup_ListMembers(t *testing.T) {
	type fields struct {
		Name    string
		Members []*models.User
	}
	tests := []struct {
		name   string
		fields fields
		want   []models.User
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				Name:    tt.fields.Name,
				Members: tt.fields.Members,
			}
			if got := g.ListMembers(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListMembers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_PrintGroupInfo(t *testing.T) {
	type fields struct {
		Name    string
		Members []*models.User
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				Name:    tt.fields.Name,
				Members: tt.fields.Members,
			}
			g.PrintGroupInfo()
		})
	}
}

func TestGroup_RemoveMember(t *testing.T) {
	type fields struct {
		Name    string
		Members []*models.User
	}
	type args struct {
		userID int32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				Name:    tt.fields.Name,
				Members: tt.fields.Members,
			}
			if err := g.RemoveMember(tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("RemoveMember() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
