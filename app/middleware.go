package app

import (
	"github.com/NikhilSharmaWe/lib/models"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
)

func (app *Application) IsAuthorized(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := []byte(app.SecretKey)
		if c.Request().Header["Token"] != nil {
			token, err := jwt.Parse(c.Request().Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return key, models.ErrUnauthorized
				}
				return key, nil
			})
			if err != nil {
				c.Logger().Error(err)
				return unauthorizedError(c)
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

				username := claims["username"].(string)
				t := claims["type"].(string)
				authorized := claims["authorized"].(bool)

				if !authorized {
					c.Logger().Error(err)
					return unauthorizedError(c)
				}

				if t != app.Users[username].Type {
					c.Logger().Error(err)
					return unauthorizedError(c)
				}

				c.Set("username", username)

				return next(c)
			} else {
				c.Logger().Error("invalid token")
				return unauthorizedError(c)
			}
		} else {
			c.Logger().Error("token header not set")
			return unauthorizedError(c)
		}
	}
}
