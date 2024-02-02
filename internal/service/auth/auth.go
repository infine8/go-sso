package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/lib/jwt"
	"sso/lib/logger/sl"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
    ErrInvalidCredentials = errors.New("invalid credentials")
)


type UserSaver interface {
    SaveUser(
        ctx context.Context,
        email string,
        passHash []byte,
    ) (uid int64, err error)
}

type UserProvider interface {  
    User(ctx context.Context, email string) (models.User, error)
    IsUserInRole(ctx context.Context, userId int64, role string) (bool, error)
}

type AppProvider interface {  
    App(ctx context.Context, appID int) (models.App, error)  
}

type Auth struct {
    log         *slog.Logger
    usrSaver    UserSaver
    usrProvider UserProvider
    appProvider AppProvider
    tokenTTL    time.Duration
}

func New(  
    log *slog.Logger,  
    userSaver UserSaver,  
    userProvider UserProvider,  
    appProvider AppProvider,  
    tokenTTL time.Duration,  
) *Auth {  
    return &Auth{  
       usrSaver:    userSaver,  
       usrProvider: userProvider,  
       log:         log,  
       appProvider: appProvider,  
       tokenTTL:    tokenTTL,
    }  
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, pass string) (int64, error) {
    const op = "Auth.RegisterNewUser"

    log := a.log.With(
        slog.String("op", op),
        slog.String("email", email),
    )

    log.Info("registering user")

    passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
    if err != nil {
        log.Error("failed to generate password hash", sl.Err(err))

        return 0, fmt.Errorf("%s: %w", op, err)
    }

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
    if err != nil {
        log.Error("failed to save user", sl.Err(err))

        return 0, fmt.Errorf("%s: %w", op, err)
    }

	return id, nil
}


func (a *Auth) Login(
    ctx context.Context,
    email string,
    password string,
    appID int,
) (string, error) {
    const op = "Auth.Login"

    log := a.log.With(
        slog.String("op", op),
        slog.String("username", email),
    )

    log.Info("attempting to login user")

	user, err := a.usrProvider.User(ctx, email)
    if err != nil {
        log.Error("failed to fetch the user", sl.Err(err))

        return "", fmt.Errorf("%s: %w", op, err)
    }

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
        log.Error("incorrect password", ErrInvalidCredentials)

        return "", fmt.Errorf("%s: %w", op, err)		
	}

	app, err := a.appProvider.App(ctx, appID)
    if err != nil {
        log.Error("failed to fetch the app", sl.Err(err))

        return "", fmt.Errorf("%s: %w", op, err)
    }

	log.Info("user logged in successfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
        log.Error("failed to create token", sl.Err(err))

        return "", fmt.Errorf("%s: %w", op, err)
    }

	return token, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userId int64) (isAdmin bool, err error) {
    const op = "Auth.IsAdmin"

    log := a.log.With(
        slog.String("op", op),
        slog.String("userId", string(userId)),
    )

    isAdmin, err = a.usrProvider.IsUserInRole(ctx, userId, "admin")
    if err != nil {
        log.Error("failed to fetch admin role", sl.Err(err))
    }

    return isAdmin, nil
}