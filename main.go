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
	mu      sync.Mutex
}

type Transaction struct {
	FromID string
	ToID   string
	Amount float64
}

type PaymentSystem struct {
	Users            map[string]*User
	TransactionQueue []Transaction
}

func (p *PaymentSystem) AddUser(u *User) {
	p.Users[u.ID] = u
}

func (p *PaymentSystem) AddTransaction(t Transaction) {
	p.TransactionQueue = append(p.TransactionQueue, t)
}

func (p *PaymentSystem) ProcessingTransactions(t Transaction) error {
	fromUser, ok := p.Users[t.FromID]
	if ok != true {
		return errors.New("fromUser not found")
	}
	toUser, ok := p.Users[t.ToID]
	if ok != true {
		return errors.New("toUser not found")
	}
	err := fromUser.Withdraw(t.Amount)
	if err != nil {
		return err
	}
	toUser.Deposit(t.Amount)
	return nil
}

func (u *User) Deposit(money float64) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.Balance = u.Balance + money
}

func (u *User) Withdraw(money float64) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.Balance < money {
		return errors.New("Error: Not enought money on balance")
	} else {
		u.Balance = u.Balance - money
		return nil
	}
}

func (p *PaymentSystem) Worker(ch <-chan Transaction, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range ch {
		err := p.ProcessingTransactions(t)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	ps := &PaymentSystem{
		Users: make(map[string]*User),
	}

	fmt.Println("Создаю UserID: 1 с балансом 1000")
	fmt.Println("Создаю UserID: 2 с балансом 500")

	user1 := &User{ID: "1", Name: "User One", Balance: 1000}
	user2 := &User{ID: "2", Name: "User Two", Balance: 500}

	ps.AddUser(user1)
	ps.AddUser(user2)

	fmt.Println("Перевожу с UserID: 1 на UserID: 2 сумму в размере 200")
	fmt.Println("Перевожу с UserID: 2 на UserID: 1 сумму в размере 50")

	tx1 := Transaction{FromID: "1", ToID: "2", Amount: 200}
	tx2 := Transaction{FromID: "2", ToID: "1", Amount: 50}

	ps.AddTransaction(tx1)
	ps.AddTransaction(tx2)

	ch := make(chan Transaction, len(ps.TransactionQueue))
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go ps.Worker(ch, &wg)
	}

	for _, tx := range ps.TransactionQueue {
		ch <- tx
	}
	close(ch)
	wg.Wait()
	fmt.Println("Итого")
	fmt.Printf("У первого пользователя должно получиться 850, а получилось: %v\n", ps.Users["1"].Balance)
	fmt.Printf("У второго пользователя должно получиться 650, а получилось: %v\n", ps.Users["2"].Balance)
}
