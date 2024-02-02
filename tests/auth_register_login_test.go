package tests

import (
	"fmt"
	"sso/tests/suite"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	ssov1 "github.com/infine8/go-sso-proto/gen/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/brianvoe/gofakeit"
)

const (
    appID = 1                 
    appSecret = "test-secret"
	passLen = 5
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

    email := gofakeit.Email()
    pass := randomFakePassword()

    // Сначала зарегистрируем нового пользователя, которого будем логинить
    respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
        Email:    email,
        Password: pass,
    })

    // Это вспомогательный запрос, поэтому делаем лишь минимальные проверки
    require.NoError(t, err)
    assert.NotEmpty(t, respReg.GetUserId())

    // А это основная проверка
    respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
        Email:    email,
        Password: pass,
        AppId:    appID,
    })
    require.NoError(t, err)

	// Получаем токен из ответа
    token := respLogin.GetToken()
    fmt.Println("token", token)
    require.NotEmpty(t, token) // Проверяем, что он не пустой

    // Отмечаем время, в которое бы выполнен логин.
    // Это понадобится для проверки TTL токена
    loginTime := time.Now()

    // Парсим и валидируем токен
    tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
        return []byte(appSecret), nil
    })
    // Если ключ окажется невалидным, мы получим соответствующую ошибку
    require.NoError(t, err)

    // Преобразуем к типу jwt.MapClaims, в котором мы сохраняли данные
    claims, ok := tokenParsed.Claims.(jwt.MapClaims)
    require.True(t, ok)

    // Проверяем содержимое токена
    assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
    assert.Equal(t, email, claims["email"].(string))
    assert.Equal(t, appID, int(claims["app_id"].(float64)))

    const deltaSeconds = 1

    logoutTime := loginTime.Add(st.Cfg.TokenTTL).Unix()
    realLogoutTime := claims["exp"].(float64)

    // Проверяем, что TTL токена примерно соответствует нашим ожиданиям.
    assert.InDelta(t, logoutTime, realLogoutTime, deltaSeconds)
}

func randomFakePassword() string {
    return gofakeit.Password(true, true, true, true, false, passLen)
}