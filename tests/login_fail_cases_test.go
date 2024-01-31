package tests

import (
	"sso/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"

	ssov1 "github.com/infine8/go-sso-proto/gen/go/sso"
)

func TestLogin_FailCases(t *testing.T) {
    ctx, st := suite.New(t)

    tests := []struct {
        name        string
        email       string
        password    string
        appID       int32
        expectedErr string
    }{
        {
            name:        "Login with Empty Password",
            email:       gofakeit.Email(),
            password:    "",
            appID:       appID,
            expectedErr: "password is required",
        },
        {
            name:        "Login with Empty Email",
            email:       "",
            password:    randomFakePassword(),
            appID:       appID,
            expectedErr: "email is required",
        },
        {
            name:        "Login with Both Empty Email and Password",
            email:       "",
            password:    "",
            appID:       appID,
            expectedErr: "email is required",
        },
        {
            name:        "Login with Non-Matching Password",
            email:       gofakeit.Email(),
            password:    randomFakePassword(),
            appID:       appID,
            expectedErr: "failed to login",
        },
        {
            name:        "Login without AppID",
            email:       gofakeit.Email(),
            password:    randomFakePassword(),
            appID:       0,
            expectedErr: "app_id is required",
        },
    }

    for _, tt := range tests {

        t.Run(tt.name, func(t *testing.T) {

            _, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
                Email:    gofakeit.Email(),
                Password: randomFakePassword(),
            })
            require.NoError(t, err)

            _, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
                Email:    tt.email,
                Password: tt.password,
                AppId:    tt.appID,
            })
            require.Error(t, err)
            require.Contains(t, err.Error(), tt.expectedErr)

        })

    }
}