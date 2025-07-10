package hasher

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"auth_service/internal/hash"
	"auth_service/pkg/conv"
	"auth_service/pkg/log"
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

func (h *hasher) CheckHash(hash, pwd string) error {
	log.Debug("CheckHash called",
		"hash", hash,
		"pwd", pwd,
		"hash_len", len(hash),
		"pwd_len", len(pwd),
		"hash_bytes_hex", fmt.Sprintf("% x", []byte(hash)),
		"pwd_bytes_hex", fmt.Sprintf("% x", []byte(pwd)),
	)

	hashBytes := conv.StrToBytes(hash)
	pwdBytes := conv.StrToBytes(pwd)

	log.Debug("Before bcrypt.CompareHashAndPassword",
		"hash_bytes_len", len(hashBytes),
		"pwd_bytes_len", len(pwdBytes),
	)

	err := bcrypt.CompareHashAndPassword(hashBytes, pwdBytes)

	if err != nil {
		log.Error("bcrypt.CompareHashAndPassword error", "error", err)
	} else {
		log.Debug("password match successful")
	}

	log.Debug("CheckHash finished",
		"error", err,
	)

	return err
}
