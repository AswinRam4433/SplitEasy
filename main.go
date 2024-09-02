package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"splitwise/group"
	"splitwise/models"
	"strconv"
	"strings"
)

//func main() {
//	fmt.Println("Hello")
//	fmt.Println("New Line")
//	user1 := models.NewUser("Aswin")
//	fmt.Println(user1)
//	user2 := models.NewUser("Bhuvan")
//	fmt.Println(user2)
//	user3 := models.NewUser("Chetan")
//	fmt.Println(user3)
//}

var (
	users    []*models.User
	groups   []*group.Group
	payments []*models.Payment
)

var (
	logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
)

func main() {
	e := echo.New()

	// Routes
	e.POST("/users", createUser)
	e.GET("/users/:id", getUser)
	e.POST("/groups", createGroup)
	e.GET("/groups/:name", getGroup)
	e.POST("/payments", createPayment)
	e.GET("/payments/:id", getPayment)
	e.GET("/list", listUsers)

	// Start server
	logger.Println("Attempting To Start Server...")
	e.Logger.Fatal(e.Start(":8080"))

}

func createUser(c echo.Context) error {
	name := c.FormValue("name")
	fmt.Println("Name is ", name)
	user := models.NewUser(name)
	users = append(users, user)
	logger.Println("Created User With Id: ", user.Id)
	return c.JSON(http.StatusCreated, user)
}

func getUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Println("Invalid ID format")
		return c.JSON(http.StatusBadRequest, "Invalid ID format")
	}

	for _, user := range users {
		if user.Id == int32(id) {
			logger.Println("Retrieved User With Id: ", user.Id)
			return c.JSON(http.StatusOK, user)
		}
	}

	logger.Println("No Matching User")
	return c.JSON(http.StatusNotFound, "User not found")
}

func createGroup(c echo.Context) error {
	name := c.FormValue("name")
	members := parseUserIDs(c.FormValue("members"))
	createdGroup := group.NewGroup(name, members)
	groups = append(groups, createdGroup)
	logger.Println("Created Group With Name: ", createdGroup.Name)
	return c.JSON(http.StatusCreated, groups)
}

func getGroup(c echo.Context) error {
	name := c.Param("name")
	for _, eachGroup := range groups {
		if eachGroup.Name == name {
			logger.Println("Retrieved Group With Name: ", eachGroup.Name)
			return c.JSON(http.StatusOK, eachGroup)

		}
	}
	logger.Println("No Matching Group")
	return c.JSON(http.StatusNotFound, "Group not found")
}

func createPayment(c echo.Context) error {
	payerID := c.FormValue("payer")
	payeeID := c.FormValue("payee")
	amountStr := c.FormValue("amount")
	mode := models.PaymentMode(c.FormValue("mode"))
	identifier := c.FormValue("identifier")
	note := c.FormValue("note")
	expenses := parseExpenseIDs(c.FormValue("expenses"))

	// Convert amount from string to float64
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		logger.Println("Invalid Amount Format")
		return c.JSON(http.StatusBadRequest, "Invalid amount format")
	}

	payerIdConv, err := strconv.ParseInt(payerID, 10, 32)
	if err != nil {
		log.Fatal("Error In Payer ID Conversion Process")
	}
	payeeIdConv, err := strconv.ParseInt(payeeID, 10, 32)
	if err != nil {
		log.Fatal("Error In Payee ID Conversion Process")
	}
	payer := findUserByID(int32(payerIdConv))
	payee := findUserByID(int32(payeeIdConv))

	if payer == nil || payee == nil {
		logger.Println("Invalid Payer or Payee")
		return c.JSON(http.StatusBadRequest, "Invalid payer or payee")
	}

	payment := models.NewPayment(payer, payee, amount, mode, identifier, note, expenses)
	payments = append(payments, payment)

	err = payment.SettlePayment()
	if err != nil {
		logger.Println("Invalid Settlement")
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	logger.Println("Created Payment")
	return c.JSON(http.StatusCreated, payment)
}

func getPayment(c echo.Context) error {
	id := c.Param("id")
	for _, payment := range payments {
		if string(payment.ID) == id {
			logger.Println("Payment Retrieved With Id: ", payment.ID)
			return c.JSON(http.StatusOK, payment)
		}
	}
	logger.Println("Payment Not Found")
	return c.JSON(http.StatusNotFound, "Payment not found")
}

// Helper functions
func parseUserIDs(userIDs string) []*models.User {
	ids := strings.Split(userIDs, ",")
	var users []*models.User

	for _, idStr := range ids {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue // Skip any IDs that cannot be converted to integers
		}

		user := findUserByID(int32(id))
		if user != nil {
			users = append(users, user)
		}
	}

	return users
}

func parseExpenseIDs(expenseIDs string) []*models.Expense {
	ids := strings.Split(expenseIDs, ",")
	var expenses []*models.Expense

	for _, idStr := range ids {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue // Skip any IDs that cannot be converted to integers
		}

		expense := findExpenseByID(int32(id))
		if expense != nil {
			expenses = append(expenses, expense)
		}
	}

	return expenses
}

func findUserByID(id int32) *models.User {
	for _, user := range users {
		if user.Id == id {
			return user
		}
	}
	return nil
}

var expenses []*models.Expense

func findExpenseByID(id int32) *models.Expense {
	for _, expense := range expenses {
		if expense.Id == id {
			return expense
		}
	}
	return nil
}

func listUsers(c echo.Context) error {
	logger.Println("Listing Users")
	for _, user := range users {
		logger.Println("User:", user.Id, user.Name)
	}
	return c.JSON(http.StatusOK, users)
}
