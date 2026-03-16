package helpers

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CompareHashPassword(hashedPwd string, rawPwd string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(rawPwd))
	if err != nil {
		return err
	}

	return nil
}
