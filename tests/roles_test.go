package tests

import (
	"sso/tests/suite"
	"testing"

	ssov1 "github.com/infine8/go-sso-proto/gen/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const(
	adminUserId = 14
)

func Test_Roles(t *testing.T) {
	ctx, st := suite.New(t)

	resp, err := st.AuthClient.IsAdmin(ctx, &ssov1.IsAdminRequest{
		UserId: adminUserId,
	})

	require.NoError(t, err)

	assert.Equal(t, true, resp.IsAdmin)
}