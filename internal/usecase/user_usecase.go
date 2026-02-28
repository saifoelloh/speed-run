package usecase

import (
	"context"
	"errors"
	"time"

	"perpustakaan/internal/domain"
	"perpustakaan/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo   domain.UserRepository
	tokenMaker jwt.TokenMaker
	ctxTimeout time.Duration
}

// NewUserUsecase creates a new instance of UserUsecase
func NewUserUsecase(ur domain.UserRepository, maker jwt.TokenMaker, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepo:   ur,
		tokenMaker: maker,
		ctxTimeout: timeout,
	}
}

func (u *userUsecase) Register(c context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	existingUser, _ := u.userRepo.GetByEmail(ctx, user.Email)
	if existingUser != nil {
		return errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return u.userRepo.Create(ctx, user)
}

func (u *userUsecase) Login(c context.Context, email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	return u.tokenMaker.CreateToken(user.ID)
}

func (u *userUsecase) GetProfile(c context.Context, id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	return u.userRepo.GetByID(ctx, id)
}
