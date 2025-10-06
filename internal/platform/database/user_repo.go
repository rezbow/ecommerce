package database

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rezbow/ecommerce/internal/app/models"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

func (repo *UserRepo) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := repo.DB.First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, ErrInternal
	}
	return &user, nil
}

func (repo *UserRepo) Get(id string) (*models.User, error) {
	var user models.User
	if err := repo.DB.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, ErrInternal
	}
	return &user, nil
}

func (repo *UserRepo) Create(user *models.User) error {
	user.ID = uuid.New()
	if err := repo.DB.Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrDuplicateKey
		}
		return ErrInternal
	}
	return nil
}
