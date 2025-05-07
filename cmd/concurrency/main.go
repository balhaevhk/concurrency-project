package main

import (
	"concurrency/internal/user"
	"fmt"
)



func main() {
	user1 := &user.User{ID: "1", Name: "John", Balance: 1000}
	user2 := &user.User{ID: "1", Name: "Petr", Balance: 3000}

	user1.Deposit(4000)
	user1.Withdraw(500)
	fmt.Println(user1)

	user2.Deposit(4000)
	err := user2.Withdraw(10000)
	if err != nil {
		fmt.Println("error:", err)
	}
	err = user2.Withdraw(5000)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(user2)
}