# SplitEasy - An Expense Management System

## Introduction

This project is a robust expense management system designed to facilitate the tracking and settling of expenses among users. It is built using the Echo framework, known for its high performance and efficiency in handling HTTP requests. The system leverages MongoDB for persistent storage and utilizes Test-Driven Development (TDD) principles to ensure code quality and reliability.

### Key Features

- **Expense Tracking:** Manage expenses by recording who paid, who shares the expense, and the remaining amount to be settled.
- **Payment Handling:** Record payments made to settle expenses, including various payment modes such as Cash, Bank Transfer, and UPI.
- **Group Management:** Organize users into groups to simplify the management of group expenses.
- **API Testing:** Endpoints have been thoroughly tested using Postman to ensure correctness and reliability.
- **Issues Tracking:** Issues encountered during development have been added and tagged for ease of development.

### Development and Testing

- **Test-Driven Development (TDD):** The application has been developed using TDD principles. Comprehensive test cases are included in the `_test.go` files for all components in the models directory.
- **GitHub Workflows:** Continuous Integration and Continuous Deployment (CI/CD) are managed through GitHub Workflows, ensuring that code changes are automatically tested and deployed.
- **Echo Framework:** The application utilizes the Echo framework for handling HTTP requests and responses, providing a fast and efficient web server.

## Entity-Relationship Diagram (ERD) Overview

### Entities

#### User

- **Attributes:**
  - `Name` (string): The name of the user.
  - `Balance` (float64): The current balance of the user.
  - `Id` (int32): A unique identifier for the user.

- **Relationships:**
  - A User can be a payer or a payee in a Payment.
  - A User can be involved in multiple Expense instances as part of the SplitBetween.

#### Expense

- **Attributes:**
  - `ID` (int): Unique identifier for the expense.
  - `Amount` (float64): The total amount of the expense.
  - `PaidBy` (*User): The user who paid for the expense.
  - `SplitBetween` ([]*User): List of users who share the expense.
  - `SplitRate` ([]float32): Rates to split the expense among users.
  - `RemainingAmount` (float64): The amount left to be settled.
  - `Payments` ([]*Payment): List of payments made towards this expense.
  - `Timestamp` (time.Time): The time when the expense was created.

- **Relationships:**
  - An Expense can be associated with multiple Payments (one-to-many).
  - An Expense involves a PaidBy user and multiple SplitBetween users.

#### Payment

- **Attributes:**
  - `ID` (int): Unique identifier for the payment.
  - `Payer` (*User): The user who made the payment.
  - `Payee` (*User): The user who received the payment.
  - `Amount` (float64): The amount paid.
  - `Mode` (PaymentMode): The mode of payment (Cash, BankTransfer, UPI).
  - `Timestamp` (time.Time): The time when the payment was made.
  - `Identifier` (string): A unique identifier for the payment.
  - `Note` (string): Additional notes for the payment.
  - `Expenses` ([]*Expense): List of expenses covered by this payment.

- **Relationships:**
  - A Payment can cover multiple Expenses (one-to-many).
  - A Payment involves a Payer and a Payee (both are users).

#### Group

- **Attributes:**
  - `Name` (string): The name of the group.
  - `Members` ([]*User): List of users in the group.
  - `Expenses` ([]*Expense): List of expenses associated with the group.

- **Relationships:**
  - A Group has multiple Members (one-to-many).
  - A Group can have multiple Expenses (one-to-many).

### Relationships and Associations

- **User ↔ Expense**
  - A User can be associated with multiple Expenses through the PaidBy and SplitBetween fields.
  - An Expense has one PaidBy user and multiple SplitBetween users.

- **User ↔ Payment**
  - A User can be a Payer or Payee in multiple Payments.
  - A Payment involves one Payer and one Payee.

- **Expense ↔ Payment**
  - An Expense can be settled by multiple Payments.
  - A Payment can cover multiple Expenses.

- **Group ↔ User**
  - A Group consists of multiple Members who are Users.

- **Group ↔ Expense**
  - A Group can have multiple Expenses associated with it.
