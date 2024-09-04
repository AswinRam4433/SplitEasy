package main

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"reflect"
	"splitwise/group"
	"splitwise/models"
	"strconv"
	"strings"
)

var (
	users    []*models.User
	groups   []*group.Group
	payments []*models.Payment
)

var (
	infoLogger  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	warnLogger  = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

var expensesMap = make(map[int]*models.Expense) // map to hold expenses by ID
var paymentsMap = make(map[int]*models.Payment) // map to hold payments by ID

func main() {
	e := echo.New()

	// Routes
	e.POST("/users", createUser)
	e.GET("/users/:id", getUser)
	e.GET("/list", listUsers)
	e.POST("/groups", createGroup)
	e.GET("/groups/:name", getGroup)
	e.POST("/payments", createPayment)
	e.GET("/payments/:id", getPayment)
	e.POST("/groups/:name/expenses", createExpense)
	e.GET("/expenses", listExpenses)

	// Start server
	infoLogger.Println("Attempting To Start Server...")
	e.Logger.Fatal(e.Start(":8080"))

}

func createUser(c echo.Context) error {
	name := c.FormValue("name")
	fmt.Println("Name is ", name)
	user := models.NewUser(name)
	users = append(users, user)
	infoLogger.Println("Created User With Id: ", user.Id)
	return c.JSON(http.StatusCreated, user)
}

func getUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		warnLogger.Println("Invalid ID format")
		return c.JSON(http.StatusBadRequest, "Invalid ID format")
	}

	for _, user := range users {
		if user.Id == int32(id) {
			infoLogger.Println("Retrieved User With Id: ", user.Id)
			return c.JSON(http.StatusOK, user)
		}
	}

	errorLogger.Println("No Matching User")
	return c.JSON(http.StatusNotFound, "User not found")
}

func createGroup(c echo.Context) error {
	name := c.FormValue("name")
	members, err := parseUserIDs(c.FormValue("members"))
	if err != nil {
		errorLogger.Println(err)
	}
	createdGroup := group.NewGroup(name, members)
	groups = append(groups, createdGroup)
	infoLogger.Println("Created Group With Name: ", createdGroup.Name)
	return c.JSON(http.StatusCreated, groups)
}

func getGroup(c echo.Context) error {
	name := c.Param("name")
	for _, eachGroup := range groups {
		if eachGroup.Name == name {
			infoLogger.Println("Retrieved Group With Name: ", eachGroup.Name)
			return c.JSON(http.StatusOK, eachGroup)

		}
	}
	errorLogger.Println("No Matching Group")
	return c.JSON(http.StatusNotFound, "Group not found")
}

func createPayment(c echo.Context) error {
	payerID := c.FormValue("payer")
	payeeID := c.FormValue("payee")
	amountStr := c.FormValue("amount")
	mode := models.PaymentMode(c.FormValue("mode"))
	identifier := c.FormValue("identifier")
	note := c.FormValue("note")
	expenseIDs := c.FormValue("expenses")

	// Convert amount from string to float64
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		warnLogger.Println("Invalid Amount Format")
		return c.JSON(http.StatusBadRequest, "Invalid amount format")
	}

	// Convert payer and payee IDs to integers
	payerIdConv, err := strconv.ParseInt(payerID, 10, 32)
	if err != nil {
		errorLogger.Println("Error In Payer ID Conversion Process")
		return c.JSON(http.StatusBadRequest, "Invalid payer ID format")
	}

	payeeIdConv, err := strconv.ParseInt(payeeID, 10, 32)
	if err != nil {
		errorLogger.Println("Error In Payee ID Conversion Process")
		return c.JSON(http.StatusBadRequest, "Invalid payee ID format")
	}

	payer := findUserByID(int32(payerIdConv))
	payee := findUserByID(int32(payeeIdConv))

	if payer == nil || payee == nil {
		warnLogger.Println("Invalid Payer or Payee")
		return c.JSON(http.StatusBadRequest, "Invalid payer or payee")
	}

	// Find expenses associated with the payment
	expenses := parseExpenseIDs(expenseIDs)

	if len(expenses) == 0 {
		warnLogger.Println("No valid expenses found")
		return c.JSON(http.StatusBadRequest, "No valid expenses found")
	}

	// Create the payment
	payment := models.NewPayment(payer, payee, amount, mode, identifier, note, expenses)
	payments = append(payments, payment)

	// Attempt to settle the payment against the expenses
	err = payment.SettlePayment()
	if err != nil {
		errorLogger.Println("Invalid Settlement:", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	infoLogger.Println("Created Payment")
	return c.JSON(http.StatusCreated, payment)
}

func getPayment(c echo.Context) error {
	id := c.Param("id")
	for _, payment := range payments {
		if strconv.Itoa(payment.ID) == id {
			infoLogger.Println("Payment Retrieved With Id: ", payment.ID)
			return c.JSON(http.StatusOK, payment)
		}
	}
	errorLogger.Println("Payment Not Found")
	return c.JSON(http.StatusNotFound, "Payment not found")
}

// Helper functions
func parseUserIDs(userIDs string) ([]*models.User, error) {
	ids := strings.Split(userIDs, ",")
	var users []*models.User
	var userAvailabilityErr error
	for _, idStr := range ids {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue // Skip any IDs that cannot be converted to integers
		}

		user := findUserByID(int32(id))

		if user != nil {
			users = append(users, user)
		} else {
			userAvailabilityErr = errors.New(fmt.Sprintf("User with ID %d not found", id))
		}
	}

	return users, userAvailabilityErr
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
		if expense.ID == int(id) {
			return expense
		}
	}
	return nil
}

func listUsers(c echo.Context) error {
	infoLogger.Println("Listing Users")
	for _, user := range users {
		infoLogger.Println("User:", user.Id, user.Name)
	}
	return c.JSON(http.StatusOK, users)
}

func createExpense(c echo.Context) error {
	fmt.Println("Received Headers: ", c.Request().Header)
	fmt.Println("Received Form Data: ", c.Request().PostForm)
	fmt.Println("Received Raw Body: ", c.Request().Body)

	amountStr := c.FormValue("amount")
	fmt.Println("Received Amount: ", amountStr, " Type: ", reflect.TypeOf(amountStr))

	if amountStr == "" {
		warnLogger.Println("Amount is missing in the form data")
		return c.JSON(http.StatusBadRequest, "Amount is missing in the form data")
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		errorLogger.Println("Invalid amount format")
		return c.JSON(http.StatusBadRequest, "Invalid amount format")
	}
	fmt.Println("Received Amount: ", c.FormValue("amount"), " ", reflect.TypeOf(c.FormValue("amount")))
	groupName := c.Param("name")

	paidByID := c.FormValue("paidBy")
	splitBetweenIDs := c.FormValue("splitBetween")
	splitRatesStr := c.FormValue("splitRates")

	// Find the payer by ID
	paidByIdConv, err := strconv.ParseInt(paidByID, 10, 32)
	if err != nil {
		errorLogger.Println(fmt.Sprint("Error in PaidBy ID conversion: ", err))
		return c.JSON(http.StatusBadRequest, "Invalid paidBy ID format")
	}

	paidBy := findUserByID(int32(paidByIdConv))
	if paidBy == nil {
		errorLogger.Println("PaidBy user not found")
		return c.JSON(http.StatusNotFound, "PaidBy user not found")
	}

	// Parse splitBetween user IDs
	splitBetweenUsers, err := parseUserIDs(splitBetweenIDs)
	if err != nil {
		errorLogger.Println(err)
	}
	if len(splitBetweenUsers) == 0 {
		errorLogger.Println("No valid users found in splitBetween")
		return c.JSON(http.StatusBadRequest, "No valid users found in splitBetween")
	}
	infoLogger.Println("Split Between Users: ", splitBetweenUsers)

	// Parse split rates
	splitRates := parseFloat32Array(splitRatesStr)
	if len(splitRates) == 0 {
		warnLogger.Println("No valid splits found in splitRates")
		return c.JSON(http.StatusBadRequest, "No valid users found in splitBetween")
	}
	infoLogger.Println("Split Rates: ", splitRates)

	if len(splitRates) != len(splitBetweenUsers) {
		errorLogger.Println("Split rates count does not match the number of users")
		return c.JSON(http.StatusBadRequest, "Invalid split rates")
	}

	// Find the group
	var group *group.Group
	for _, g := range groups {
		if g.Name == groupName {
			group = g
			break
		}
	}

	if group == nil {
		warnLogger.Println("Group not found")
		return c.JSON(http.StatusNotFound, "Group not found")
	}

	// Create the expense
	expense := models.NewExpense(amount, paidBy, splitBetweenUsers, splitRates)

	group.AddExpense(expense)
	expenses = append(expenses, expense)
	expensesMap[expense.ID] = expense

	// Split the expense to update the balances
	err = expense.SplitExpense()
	if err != nil {
		errorLogger.Println("Error splitting expense in CreateExpense:", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	infoLogger.Println("Added Expense to Group:", group.Name)
	return c.JSON(http.StatusCreated, expense)
}

// Helper function to parse comma-separated float32 values
func parseFloat32Array(input string) []float32 {
	strValues := strings.Split(input, ",")
	var floatValues []float32

	for _, str := range strValues {
		if value, err := strconv.ParseFloat(str, 32); err == nil {
			floatValues = append(floatValues, float32(value))
		}
	}

	return floatValues
}

func listExpenses(c echo.Context) error {
	infoLogger.Println("Listing Expenses")
	for _, expense := range expenses {
		infoLogger.Println("Expense ID:", expense.ID, "Amount:", expense.Amount)
	}
	return c.JSON(http.StatusOK, expenses)
}
