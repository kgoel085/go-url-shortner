package utils

import "golang.org/x/crypto/bcrypt"

func HashPwd(pwd string) (string, error) {
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashPwd), nil
}

func CheckHashPwd(pwd, hashPwd string) bool {
	hashPwdErr := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(pwd))
	return hashPwdErr == nil
}
