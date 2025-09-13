package userrepo

import (
	"context"

	errordom "github.com/bsach64/booked/internal/domain/error"
	userdom "github.com/bsach64/booked/internal/domain/user"
	"github.com/bsach64/booked/internal/repo/sql/db"
	"github.com/bsach64/booked/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type impl struct {
	config *utils.Config
	dbConn *db.Queries
}

func (i *impl) CreateUser(ctx context.Context, user userdom.User) error {
	createUserParams := db.CreateUserParams{
		ID:             pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Name:           user.Name,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
	}

	if user.Role != userdom.USER && user.Role != userdom.ADMIN {
		return errordom.GetUserError(errordom.INVALID_USER_ROLE, "not a valid user role", nil)
	}
	createUserParams.Role = db.UserRole(user.Role)

	err := i.dbConn.CreateUser(ctx, createUserParams)
	if err != nil {
		return errordom.GetDBError(errordom.DB_WRITE_ERROR, "could not create user", err)
	}
	return nil
}

func (i *impl) GetUserByEmail(ctx context.Context, email string) (*userdom.User, error) {
	dbUser, err := i.dbConn.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, errordom.GetDBError(errordom.DB_READ_ERROR, "could not get user from email", err)
	}
	return ToUserDomain(dbUser), nil
}

func New(config *utils.Config, db *db.Queries) userdom.Repository {
	return &impl{
		config: config,
		dbConn: db,
	}
}
