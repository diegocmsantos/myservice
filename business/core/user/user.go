// Package user provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/diegocmsantos/myservice/business/data/store/user"
	"github.com/diegocmsantos/myservice/business/sys/auth"
	"github.com/diegocmsantos/myservice/business/sys/database"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("user not found")
	ErrInvalidID             = errors.New("ID is not in its proper form")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

// Core manages the set of APIs for user access.
type Core struct {
	log  *zap.SugaredLogger
	user user.Store
}

// NewCore constructs a core for user api access.
func NewCore(log *zap.SugaredLogger, sqlxDB *sqlx.DB) Core {
	return Core{
		log:  log,
		user: user.NewStore(log, sqlxDB),
	}
}

// Create inserts a new user into the database.
func (c Core) Create(ctx context.Context, nu user.NewUser, now time.Time) (user.User, error) {

	// PERFORM PRE BUSINESS OPERATIONS

	usr, err := c.user.Create(ctx, nu, now)
	if err != nil {
		return user.User{}, fmt.Errorf("create: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return usr, nil
}

// Update replaces a user document in the database.
func (c Core) Update(ctx context.Context, claims auth.Claims, userID string, uu user.UpdateUser, now time.Time) error {

	// PERFORM PRE BUSINESS OPERATIONS

	err := c.user.Update(ctx, claims, userID, uu, now)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return nil
}

// Delete removes a user from the database.
func (c Core) Delete(ctx context.Context, claims auth.Claims, userID string) error {

	// PERFORM PRE BUSINESS OPERATIONS

	if err := c.user.Delete(ctx, claims, userID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return nil
}

// Query retrieves a list of existing users from the database.
func (c Core) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]user.User, error) {

	// PERFORM PRE BUSINESS OPERATIONS

	dbUsers, err := c.user.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return dbUsers, nil
}

// QueryByID gets the specified user from the database.
func (c Core) QueryByID(ctx context.Context, claims auth.Claims, userID string) (user.User, error) {

	// PERFORM PRE BUSINESS OPERATIONS

	dbUsr, err := c.user.QueryByID(ctx, claims, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return user.User{}, ErrNotFound
		}
		return user.User{}, fmt.Errorf("query: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return dbUsr, nil
}

// QueryByEmail gets the specified user from the database by email.
func (c Core) QueryByEmail(ctx context.Context, claims auth.Claims, email string) (user.User, error) {

	// PERFORM PRE BUSINESS OPERATIONS

	dbUsr, err := c.user.QueryByEmail(ctx, claims, email)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return user.User{}, ErrNotFound
		}
		return user.User{}, fmt.Errorf("query: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return dbUsr, nil
}

// Authenticate finds a user by their email and verifies their password. On
// success it returns a Claims User representing this user. The claims can be
// used to generate a token for future authentication.
func (c Core) Authenticate(ctx context.Context, now time.Time, email, password string) (auth.Claims, error) {

	// PERFORM PRE BUSINESS OPERATIONS

	claims, err := c.user.Authenticate(ctx, now, email, password)
	if err != nil {
		return auth.Claims{}, fmt.Errorf("query: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return claims, nil
}
