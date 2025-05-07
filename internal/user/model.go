package user

import "errors"

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