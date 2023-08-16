package utility

import (
	uuid58 "github.com/AlexanderMatveev/go-uuid-base58"
	"github.com/google/uuid"
)

func GetRandomBase58StringOfLength(length int) (string, error) {

	uid := uuid.New()

	s, err := uuid58.ToBase58(uid)

	if err != nil {
		return "", err
	}

	return s[:length], nil
}
