package tests

import (
	"sso/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"

	ssov1 "github.com/infine8/go-sso-proto/gen/go/sso"
)

func TestRegister_FailCases(t *testing.T) {
    ctx, st := suite.New(t)

    tests := []struct {
        name        string
        email       string
        password    string
        expectedErr string
    }{
        {
            name:        "Register with Empty Password",
            email:       gofakeit.Email(),
            password:    "",
            expectedErr: "password is required",
        },
        {
            name:        "Register with Empty Email",
            email:       "",
            password:    randomFakePassword(),
            expectedErr: "email is required",
        },
        {
            name:        "Register with Both Empty",
            email:       "",
            password:    "",
            expectedErr: "email is required",
        },
    }

    for _, tt := range tests {
		
        t.Run(tt.name, func(t *testing.T) {

            _, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
                Email:    tt.email,
                Password: tt.password,
            })

            require.Error(t, err)
            require.Contains(t, err.Error(), tt.expectedErr)

        })
    }
}