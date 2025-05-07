package main

import (
	"errors"
	"fmt"
)

type User struct {
	ID string
	Name string
	Balance float64
}

func (u *User) Deposit(amount float64) {
	u.Balance += amount 
}

func (u *User) Withdraw(amount float64) error {
	if u.Balance < amount {
		return errors.New("insufficient funds on balance")
	}
	u.Balance -= amount
	return nil
}

func main() {
	user1 := &User{ID: "1", Name: "John", Balance: 1000}
	user2 := &User{ID: "1", Name: "Petr", Balance: 3000}

	user1.Deposit(4000)
	fmt.Println(user1)
	user2.Deposit(4000)
	fmt.Println(user2)
}