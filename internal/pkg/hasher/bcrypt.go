package hasher

import "golang.org/x/crypto/bcrypt"

type bcryptStrategy struct{}

func NewBcrypt() Strategy {
	return &bcryptStrategy{}
}

func (s *bcryptStrategy) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *bcryptStrategy) Compare(password string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
