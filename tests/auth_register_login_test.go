package tests

import (
	authservice "authService/protos/gen/go/sso"
	"authService/tests/suite"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	emptyAppID            = 0
	appID          uint32 = 1
	appSecret             = "secret"
	passDefaultLen        = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	username := gofakeit.Username()
	pass := RandomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &authservice.RegisterRequest{
		Username: username,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, username, respReg.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &authservice.LoginRequest{
		Username: username,
		Password: pass,
		AppId:    appID,
	})
	require.NoError(t, err)

	accessToken := respLogin.GetAccessToken()
	require.NotEmpty(t, accessToken)

	loginTime := time.Now()

	tokenParsed, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, appID, uint32(claims["app_id"].(float64)))
	assert.Equal(t, username, claims["username"].(string))
	assert.Equal(t, respReg.GetUserId(), uint32(claims["uid"].(float64)))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func RandomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
