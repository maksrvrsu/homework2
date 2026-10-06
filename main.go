package main

import (
	"errors"
	"fmt"
)

type User struct {
	ID      string
	Name    string
	Balance float64
}

func (u *User) Deposit(money float64) {
	u.Balance = u.Balance + money
}

func (u *User) Withdraw(money float64) error {
	if u.Balance < money {
		return errors.New("Error: Not enought money on balance")
	} else {
		u.Balance = u.Balance - money
		return nil
	}
}

func main() {
	user1 := &User{ID: "513882", Name: "Alexey", Balance: 0}
	user2 := &User{ID: "513883", Name: "Anton", Balance: 0}
	user1.Deposit(1000)
	fmt.Println("Выполнение метода Deposit над user1", user1.Balance)
	user2.Deposit(200)
	fmt.Println("Выполнение метода Deposit над user2", user2.Balance)
	err := user1.Withdraw(500)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Выполнение метода Withdraw над user1", user1.Balance)
	err = user2.Withdraw(500)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Выполнение метода Withdraw над user2", user2.Balance)
}
