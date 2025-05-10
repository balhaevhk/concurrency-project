package main

import (
	"errors"
	"fmt"
	"sync"
)

type User struct {
	ID      string
	Name    string
	Balance float64
	mu sync.Mutex
}

func (u *User) Deposit(amount float64) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.Balance += amount
}

func (u *User) Withdraw(amount float64) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.Balance < amount {
		return errors.New("insufficient funds on balance")
	}
	u.Balance -= amount
	return nil
}

type Transaction struct {
	FromID string
	ToID   string
	Amount float64
}

type PaymentSystem struct {
	Users        map[string]*User
	Transactions []Transaction
}

func (ps *PaymentSystem) AddUser(u *User) {
	if ps.Users == nil {
		ps.Users = make(map[string]*User)
	}
	ps.Users[u.ID] = u
}

func (ps *PaymentSystem) AddTransaction(t Transaction) {
	ps.Transactions = append(ps.Transactions, t)
}

func (ps *PaymentSystem) ProcessingTransactions(t Transaction) error {
	fromUser, fromExist := ps.Users[t.FromID]
	toUser, toExist := ps.Users[t.ToID]

	if !fromExist {
		return fmt.Errorf("user with ID %s not found", t.FromID)
	}
	if !toExist {
		return fmt.Errorf("user with ID %s not found", t.ToID)
	}

	if err := fromUser.Withdraw(t.Amount); err != nil {
		return fmt.Errorf("error withdrawing from user %s: %v", t.FromID, err)
	}

	toUser.Deposit(t.Amount)

	ps.Users[t.FromID] = fromUser
	ps.Users[t.ToID] = toUser

	return nil
}

func (ps *PaymentSystem) Worker(ch <-chan Transaction, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range ch {
		err := ps.ProcessingTransactions(t)
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func main() {

	ps := &PaymentSystem{
		Users:        make(map[string]*User),
		Transactions: []Transaction{},
	}

	fmt.Println("Создаю UserID: 1 с балансом 1000")
	fmt.Println("Создаю UserID: 2 с балансом 500")
	user1 := &User{ID: "1", Name: "John", Balance: 1000}
	user2 := &User{ID: "2", Name: "Petr", Balance: 500}

	ps.AddUser(user1)
	ps.AddUser(user2)

	fmt.Println("Перевожу с UserID: 1 на UserID: 2 сумму в размере 200")
	fmt.Println("Перевожу с UserID: 2 на UserID: 1 сумму в размере 50")

	ch := make(chan Transaction, len(ps.Transactions))

	ps.AddTransaction(Transaction{"1", "2", 200})
	ps.AddTransaction(Transaction{"2", "1", 50})

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go ps.Worker(ch, &wg)
	}

	for _, t := range ps.Transactions {
		ch <- t
	}

	close(ch)
	wg.Wait()

	fmt.Println("Итого")
	fmt.Println("John's Balance:", ps.Users["1"].Balance)
	fmt.Println("Petr's Balance:", ps.Users["2"].Balance)
	fmt.Println("У первого пользователя должно получиться 850")
	fmt.Println("У второго пользователя должно получиться 650")

}
