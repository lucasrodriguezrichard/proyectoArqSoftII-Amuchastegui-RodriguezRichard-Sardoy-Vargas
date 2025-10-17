package crypto

import (
	"github.com/blassardoy/restaurant-reservas/users-api/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

// BcryptHasher implements domain.PasswordHasher using bcrypt.
type BcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost int) *BcryptHasher {
	if cost <= 0 {
		cost = bcrypt.DefaultCost
	}
	return &BcryptHasher{cost: cost}
}

func (h *BcryptHasher) Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), h.cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (h *BcryptHasher) Compare(hash string, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

// Verify interface compliance at compile time.
var _ domain.PasswordHasher = (*BcryptHasher)(nil)
