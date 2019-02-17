package common

import (
	"fmt"

	"github.com/joho/godotenv"
)

// EncryptPassword encrypts an account holders password for database storage
func EncryptPassword(pwd string) (string, error) {
	return pwd, nil
}

func main() {
	x := "simon"
	fmt.Println(EncryptPassword(x))
	// Load environment
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}
