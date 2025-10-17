package repository

import (
	"context"
	"strings"

	"github.com/blassardoy/restaurant-reservas/users-api/internal/domain"
	"gorm.io/gorm"
)

// GormUserRepository is a GORM-backed implementation of domain.UserRepository.
type GormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) GetByID(ctx context.Context, id uint64) (domain.User, error) {
	var u domain.User
	tx := r.db.WithContext(ctx).First(&u, id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, tx.Error
	}
	return u, nil
}

func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var u domain.User
	tx := r.db.WithContext(ctx).
		Where("LOWER(email) = ?", strings.ToLower(email)).
		First(&u)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, tx.Error
	}
	return u, nil
}

func (r *GormUserRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	var u domain.User
	tx := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&u)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, tx.Error
	}
	return u, nil
}

func (r *GormUserRepository) Create(ctx context.Context, u *domain.User) error {
	tx := r.db.WithContext(ctx).Create(u)
	return tx.Error
}

func (r *GormUserRepository) UpdatePasswordHash(ctx context.Context, id uint64, newHash string) error {
	tx := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", id).
		Update("password_hash", newHash)
	return tx.Error
}

func (r *GormUserRepository) ExistsByEmailOrUsername(ctx context.Context, email, username string) (bool, error) {
	var count int64
	tx := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("LOWER(email) = ? OR username = ?", strings.ToLower(email), username).
		Count(&count)
	return count > 0, tx.Error
}
