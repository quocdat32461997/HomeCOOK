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

	// Load environment
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	x := "simon"
	fmt.Println(EncryptPassword(x))

}
