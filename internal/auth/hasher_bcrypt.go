package auth

import (
	"github.com/rotisserie/eris"
	"golang.org/x/crypto/bcrypt"
)

type hasherBcrypt struct {
	cost int
}

func NewHasherBcrypt(cost int) Hasher {
	return &hasherBcrypt{cost: cost}
}

func (hb *hasherBcrypt) Hash(val string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(val), hb.cost)
	if err != nil {
		return "", eris.Wrap(err, "error hashing value")
	}

	return string(hash), nil
}

func (hb *hasherBcrypt) CheckHash(hash, val string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(val))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}

		return false, eris.Wrap(err, "error checking hash")
	}

	return true, nil
}
