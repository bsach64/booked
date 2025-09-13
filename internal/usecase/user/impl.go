package useruc

import (
	"context"
	"time"

	errordom "github.com/bsach64/booked/internal/domain/error"
	userdom "github.com/bsach64/booked/internal/domain/user"
	"github.com/bsach64/booked/internal/repo"
	"github.com/bsach64/booked/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type impl struct {
	config       *utils.Config
	repositories repo.Repositories
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func (i *impl) LoginUser(ctx context.Context, loginRequest *userdom.LoginRequest) (*userdom.LoginResponse, error) {
	// 1 day for now
	expirationTime := time.Now().Add(24 * time.Hour)

	user, err := i.repositories.User.GetUserByEmail(ctx, loginRequest.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errordom.GetUserError(errordom.USER_NOT_FOUND, "no user with given email", nil)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(loginRequest.UnhashedPassword))
	if err != nil {
		return nil, errordom.GetUserError(errordom.INVALID_PASSWORD, "given password is incorrect", nil)
	}

	claims := &userdom.Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(i.config.JwtSecret))
	if err != nil {
		return nil, errordom.GetUserError(errordom.JWT_SIGN_FAILURE, "could not sign jwt", err)
	}
	return &userdom.LoginResponse{Token: signedToken}, nil
}

func (i *impl) CreateUser(ctx context.Context, newUser userdom.NewUser) error {
	if newUser.UnhashedPassword == "" {
		return errordom.GetUserError(errordom.EMTPY_PASSWORD, "empty password provided", nil)
	}

	if newUser.Name == "" {
		return errordom.GetUserError(errordom.EMPTY_NAME, "empty name provided", nil)
	}

	if newUser.Email == "" {
		return errordom.GetUserError(errordom.EMPTY_EMAIL, "empty email provided", nil)
	}

	hashedPassword, err := HashPassword(newUser.UnhashedPassword)
	if err != nil {
		return err
	}

	user := userdom.User{
		Name:           newUser.Name,
		Email:          newUser.Email,
		HashedPassword: hashedPassword,
		Role:           userdom.USER,
	}
	return i.repositories.User.CreateUser(ctx, user)
}

func (i *impl) GetUserByEmail(ctx context.Context, email string) (*userdom.User, error) {
	user, err := i.repositories.User.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errordom.GetUserError(errordom.USER_NOT_FOUND, "no user with given email", nil)
	}
	return user, nil
}

func New(config *utils.Config, repositories repo.Repositories) userdom.Usecase {
	return &impl{
		config:       config,
		repositories: repositories,
	}
}
