package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rezbow/ecommerce/internal/app/models"
	"github.com/rezbow/ecommerce/internal/platform/authentication"
	"github.com/rezbow/ecommerce/internal/platform/database"
)

type UserSvc struct {
	userRepo  models.UserRepo
	jwtSecret string
}

func NewUserService(repo models.UserRepo, jwtSecret string) *UserSvc {
	return &UserSvc{
		userRepo:  repo,
		jwtSecret: jwtSecret,
	}
}

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrInternal           = errors.New("internal error")
	ErrInvalidCredentials = errors.New("wrong email or password")
)

func (svc *UserSvc) RegisterUser(data *models.RegisterUser) (*models.User, error) {
	passwordHash, err := authentication.HashPassword(data.Password)
	if err != nil {
		return nil, ErrInternal
	}
	user := models.User{
		Email:        data.Email,
		PasswordHash: passwordHash,
	}
	err = svc.userRepo.Create(&user)
	if err != nil {
		if errors.Is(err, database.ErrDuplicateKey) {
			return nil, ErrDuplicateEmail
		}
		return nil, ErrInternal
	}
	return &user, nil
}

func (svc *UserSvc) Authenticate(data *models.Login) (string, error) {
	// check db
	user, err := svc.userRepo.GetByEmail(data.Email)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", ErrInternal
	}

	if !authentication.CheckPassword(data.Password, user.PasswordHash) {
		return "", ErrInvalidCredentials
	}

	// generate jwt token
	token, err := authentication.NewJWTToken(user.ID, user.IsAdmin, svc.jwtSecret)
	if err != nil {
		return "", ErrInternal
	}
	return token, nil

}

func (svc *UserSvc) GetUser(id uuid.UUID) (*models.User, error) {
	user, err := svc.userRepo.Get(id.String())
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternal
	}
	return user, nil
}
