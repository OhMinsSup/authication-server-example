package lib

import "golang.org/x/crypto/bcrypt"

// JSON alias type
type JSON = map[string]interface{}

// Hash 패스워드 값을 해시값으로 리턴
func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// Check 패스워드와 hash값을 비교해서 같은 값인지 체크
func Compare(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}