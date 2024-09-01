package main

import (
	"fmt"
	"splitwise/models"
)

func main() {
	fmt.Println("Hello")
	fmt.Println("New Line")
	user1 := models.NewUser("Aswin")
	fmt.Println(user1)
	user2 := models.NewUser("Bhuvan")
	fmt.Println(user2)
	user3 := models.NewUser("Chetan")
	fmt.Println(user3)
}
