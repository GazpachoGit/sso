package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/GazpachoGit/sso/internal/domain/models"
	"github.com/GazpachoGit/sso/internal/lib/jwt"
	"github.com/GazpachoGit/sso/internal/lib/logger/sl"
	"github.com/GazpachoGit/sso/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
}

// getter service interface
type UserProvider interface {
	User(ctx context.Context, email string) (*models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

// getter service interface
type AppProvider interface {
	App(ctx context.Context, appID int64) (*models.App, error)
}

func New(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (token string, err error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	log.Info("logging in user")
	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))
			return "", fmt.Errorf("%s,%w", op, ErrInvalidCredentials)
		}
		log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password)); err != nil {
		log.Info("invalid password", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)

	}
	app, err := a.appProvider.App(ctx, int64(appID))
	if err != nil {
		log.Info("invalid app id", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("user logged in successfully")
	jwt, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("can't generate JWT", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return jwt, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error) {
	const op = "Auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		//can break some local regulatories
		slog.String("email", email),
	)
	log.Info("registering new user")
	//returns solt+hash in one string
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	userID, err = a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("user registered successfully", slog.Int64("userID", userID))
	return userID, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "Auth.IsAdmin"
	log := a.log.With(
		slog.String("op", op),
		slog.Int64("userID", userID),
	)
	log.Info("checking if user admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		log.Error("failed to check if user is admin", sl.Err(err))
		return false, err
	}
	log.Info("checked if user admin", slog.Bool("isAdmin", isAdmin))
	return isAdmin, nil
}
