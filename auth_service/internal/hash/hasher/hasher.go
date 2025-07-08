package hasher

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"auth_service/internal/hash"
	"auth_service/pkg/conv"
)

const MinHashCost = 10
const MaxHashCost = bcrypt.MaxCost

type hasher struct {
	cost int
}

func NewHasher(cost int) hash.Hasher {
	if cost < MinHashCost {
		cost = MinHashCost
	}
	if cost > MaxHashCost {
		cost = MaxHashCost
	}

	return &hasher{cost: cost}
}

func (h *hasher) Hash(pwd string) (string, error) {
	bytePwd := conv.StrToBytes(pwd)
	hashPwd, err := bcrypt.GenerateFromPassword(bytePwd, h.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return conv.BytesToStr(hashPwd), nil
}

func (h *hasher) CheckHash(pwd, hash string) error {
	return bcrypt.CompareHashAndPassword(conv.StrToBytes(hash), conv.StrToBytes(pwd))
}
