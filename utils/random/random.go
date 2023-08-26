package random

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	curr "github.com/tech_school/simple_bank/utils/currency"
	"gopkg.in/loremipsum.v1"
)

var rng *rand.Rand

func init() {
	seed := time.Now().UnixNano()
	rng = rand.New(rand.NewSource(seed))

	// Now you can use rng to generate random numbers
	randomValue := rng.Intn(100)
	println(randomValue)
}

func RandomInt(min, max int64) int64 {
	return min + rng.Int63n(max-min+1)
}

func RandomString(n int, word string) string {
	var sb strings.Builder

	lenWord := len(word)

	for i := 0; i < n; i++ {
		letter := word[rng.Intn(lenWord)]
		sb.WriteByte(letter)
	}

	return sb.String()
}

func RandomOwner() string {
	const vocal = "aiueo"
	const consonant = "bcdfghjklmnpqrstvwxyz"

	return (RandomString(1, consonant) + RandomString(1, vocal) +
		RandomString(1, consonant) + RandomString(1, vocal) +
		RandomString(1, consonant) + RandomString(1, vocal))
}

func RandomEmail(name string) string {
	return fmt.Sprintf("%s@gmail.com", name)
}

func RandomPassword() string {
	const alphabet = "abcedfghijklmnopqrstuvwxyz"
	return RandomString(4, alphabet)
}

func RandomMoney() int64 {
	return RandomInt(100, 1000)
}

func RandomCurrency() string {
	currencies := []string{curr.USD, curr.JPY, curr.RUB}
	RandomCurrencies := currencies[rng.Intn(len(currencies))]
	return RandomCurrencies
}

func RandomDescription() string {
	loremIpsumGenerator := loremipsum.New()
	return loremIpsumGenerator.Words(int(RandomInt(20, 40)))
}

// Just for fun

func RandomCountry() string {
	countries := []string{"USA", "Russia", "Japanese"}
	RandomCountries := countries[rng.Intn(len(countries))]
	return RandomCountries
}

type RdPerson struct {
	Username    string
	FullName    string
	Email       string
	Password    string
	Balance     int64
	Country     string
	Currency    string
	Description string
}

func RandomPerson() RdPerson {
	var currency string
	var lastName string
	var balance int64

	jp := []string{"Tanaka", "Hayashi", "Watanabe", "Yamada", "Kobayashi", "Nakamura", "Kimura", "Yoshida", "Shimizu", "Miyamoto"}
	us := []string{"Montgomery", "Levine", "Graham", "Lincoln", "Lopez", "Nimitz", "Ford", "Wilson", "Lennon", "Rosevelt"}
	rn := []string{"Gerasimov", "Konasenkov", "Zhukov", "Kuznetsov", "Pavlichenko", "Ovcharenko", "Sereda", "Petrov", "Prigozhin", "Gurevich"}

	country := RandomCountry()

	if country == "USA" {
		lastName = us[rng.Intn(len(us))]
		balance = RandomInt(100, 1_000)
		currency = curr.USD
	} else if country == "Russia" {
		lastName = rn[rng.Intn(len(rn))]
		balance = RandomInt(10_000, 100_000)
		currency = curr.RUB
	} else {
		lastName = jp[rng.Intn(len(jp))]
		balance = RandomInt(15_000, 150_000)
		currency = curr.JPY
	}

	name := RandomOwner()

	return RdPerson{
		Username:    name,
		FullName:    name + " " + lastName,
		Email:       RandomEmail(name),
		Password:    RandomPassword(),
		Balance:     balance,
		Country:     country,
		Currency:    currency,
		Description: RandomDescription(),
	}
}
