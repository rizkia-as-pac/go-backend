package pass

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword will compute bycript and return the bycript hash of the passsword
func HashedPassword(password string) (hashedPassword string, err error) {
	// bcrypt.GenerateFromPassword(convert password string into byte slice, cost)
	byteHashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // bcrypt.DefaultCost = 10
	if err != nil {
		return "", fmt.Errorf("gagal untuk meng hash password : %w", err)
	}
	// convert hash password dari yang berbentuk byte slice menjadi string kembali. dan kembalikan error dengan nilai nil
	hashedPassword = string(byteHashedPassword)
	return hashedPassword, nil

}

// CheckPassword check if the provided password is correct or not when compare to the provided hashed password
func CheckPassword(password string, hashPassword string) error {
	// mengambil data cost dan salt dari hashpassword yang disimpan didatabase, lalu dengan cost dan salt tersebut kita jadikan patokan untuk mengubah password yang baru dikirimkan oleh user menjadi hash baru, karena dua duanya sudah berbentuk hash maka kini bisa dikomparasikan

	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}
