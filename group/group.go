package group

import (
	"errors"
	"fmt"
	"splitwise/models"
)

type Group struct {
	Name     string
	Members  []*models.User
	Expenses []*models.Expense // To keep track of all expenses related to the group
}

func NewGroup(name string, members []*models.User) *Group {
	return &Group{
		Name:     name,
		Members:  members,
		Expenses: []*models.Expense{},
	}
}

func (g *Group) AddMember(user *models.User) {
	g.Members = append(g.Members, user)

}

func (g *Group) RemoveMember(userID int32) error {
	for i, member := range g.Members {
		if member.Id == userID {
			g.Members = append(g.Members[:i], g.Members[i+1:]...)
			// to remove member, present at index i, and to concatenate the remaining
			return nil
		}
	}
	return errors.New("member not found")
}

func (g *Group) ListMembers() []models.User {
	users := []models.User{}
	for _, member := range g.Members {
		users = append(users, *member)
	}
	return users
}

// PrintGroupInfo prints the group's name and its members
func (g *Group) PrintGroupInfo() {
	fmt.Printf("Group Name: %s\n", g.Name)
	for _, member := range g.Members {
		fmt.Printf("ID: %d, Name: %s, Balance: %.2f\n", member.Id, member.Name, member.Balance)
	}
}
