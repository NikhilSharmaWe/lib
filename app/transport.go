package app

import (
	"encoding/json"
	"net/http"

	"github.com/NikhilSharmaWe/lib/models"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func (app *Application) Router() *echo.Echo {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	e.POST("/login", app.HandleLogin)
	e.GET("/home", app.HandleHome, app.IsAuthorized)
	e.POST("/addBook", app.HandleAddBook, app.IsAuthorized)
	e.POST("/deleteBook", app.HandleDeleteBook, app.IsAuthorized)

	return e
}

func (app *Application) HandleLogin(c echo.Context) error {
	u := models.User{}
	if err := c.Bind(&u); err != nil {
		c.Logger().Error(err)
		return internalServerError(c)
	}

	user, ok := app.Users[u.Username]
	if !ok {
		return c.JSON(http.StatusBadRequest, jsonError(models.ErrUserNotFound))
	}

	if user.Password != u.Password {
		return c.JSON(http.StatusBadRequest, jsonError(models.ErrWrongPassword))
	}

	tokenString, err := app.generateJWT(user)
	if err != nil {
		c.Logger().Error(err)
		return internalServerError(c)
	}

	return c.JSON(http.StatusOK, map[string]string{"user": u.Username, "token": tokenString})
}

func (app *Application) HandleHome(c echo.Context) error {
	var (
		username = c.Get("username").(string)
		t        = app.Users[username].Type
		data     [][]string
	)

	switch t {
	case "regular":
		d, err := readCSV("regularUser.csv")
		if err != nil {
			return internalServerError(c)
		}

		data = d[1:]
	case "admin":
		da, err := readCSV("adminUser.csv")
		if err != nil {
			c.Logger().Error(err)
			return internalServerError(c)
		}

		dr, err := readCSV("regularUser.csv")
		if err != nil {
			c.Logger().Error(err)
			return internalServerError(c)
		}

		data = append(da[1:], dr[1:]...)
	default:
		return internalServerError(c)
	}

	return c.JSON(http.StatusOK, map[string]any{"Books": data})
}

func (app *Application) HandleAddBook(c echo.Context) error {
	var (
		username = c.Get("username").(string)
		t        = app.Users[username].Type
	)

	switch t {
	case "regular":
		return c.JSON(http.StatusForbidden, map[string]string{"error": models.ErrUserNotAdmin.Error()})
	case "admin":
		data := []models.Book{}
		if err := json.NewDecoder(c.Request().Body).Decode(&data); err != nil {
			c.Logger().Error(err)
			return internalServerError(c)
		}

		if err := writeBookToCSV("regularUser.csv", data); err != nil {
			c.Logger().Error(err)
			return internalServerError(c)
		}

		return nil

	default:
		return internalServerError(c)
	}
}

func (app *Application) HandleDeleteBook(c echo.Context) error {
	var (
		username = c.Get("username").(string)
		t        = app.Users[username].Type
	)

	switch t {
	case "regular":
		return c.JSON(http.StatusForbidden, map[string]string{"error": models.ErrUserNotAdmin.Error()})
	case "admin":
		data := make(map[string]interface{})

		if err := json.NewDecoder(c.Request().Body).Decode(&data); err != nil {
			c.Logger().Error(err)
			return internalServerError(c)
		}

		bookName := data["book_name"].(string)

		if err := deleteBookFromCSV("regularUser.csv", bookName); err != nil {
			if err == models.ErrBookNotExists {
				return c.JSON(http.StatusBadRequest, jsonError(models.ErrBookNotExists))
			}

			c.Logger().Error(err)
			return internalServerError(c)
		}

		return nil

	default:
		return c.JSON(http.StatusInternalServerError, models.ErrUnexpected)
	}
}
