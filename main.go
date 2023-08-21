package main

import (
	"github.com/tech_school/simple_bank/utils/random"
)

func main() {
	// fmt.Println("Hello, World!")
	// owner := random.RandomOwner()
	// println("Owner : ", owner)
	// println("Email : ", random.RandomEmail(owner))
	// println("Password : ", random.RandomPassword())
	// println("Balance : ", random.RandomMoney())
	// println("Currency : ", random.RandomCurrency())
	// println("Description : ", random.RandomDescription())

	// person := random.RandomPerson()
	// println(person.Name)
	// println(person.Email)
	// println(person.Password)
	// println(person.Balance)
	// println(person.Country)
	// println(person.Currency)
	// println(person.Description)

    for i := 0; i < 10; i++ {
        println(random.RandomOwner())
    }
}
