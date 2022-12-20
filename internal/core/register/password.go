package register

import (
	"git.friends.com/PetLand/UserService/v2/internal/models"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" + "!@#$%*?"
	length  = 60
)

func ComparePassword(hash string, password string, salt []byte) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(saltPassword(salt, password))) == nil
}

func EncryptPassword(d *models.Data) error {
	salt := createSalt()
	b, err := bcrypt.GenerateFromPassword([]byte(saltPassword(salt, d.PasswordEncoded)), bcrypt.DefaultCost)
	if err != nil {
		return err //nolint:wrapcheck
	}
	d.PasswordEncoded = string(b)
	d.PasswordSalt = string(salt)

	return nil
}

func createSalt() []byte {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	salt := make([]byte, length)
	for i := range salt {
		salt[i] = charset[seededRand.Int63()%int64(len(charset))]
	}

	return salt
}

func saltPassword(salt []byte, pwd string) string {
	b := make([]byte, len(salt))
	for i, v := range []byte(pwd) {
		b[i] = v
	}
	c := make([]byte, len(salt))
	for i := range salt {
		c[i] = salt[i] ^ b[i]
	}

	res := string(c)

	return res
}
