package utils

import "golang.org/x/crypto/bcrypt"

func HashPwd(pwd string) (string, error) {
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashPwd), nil
}
