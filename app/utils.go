package app

import (
	"encoding/csv"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/NikhilSharmaWe/lib/models"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
)

type Application struct {
	Addr      string
	SecretKey string
	Users     map[string]models.User
}

func NewApplication() *Application {
	return &Application{
		Addr:      os.Getenv("ADDR"),
		SecretKey: os.Getenv("SECRET_KEY"),
		Users: map[string]models.User{
			"Nikhil": {
				Username: "Nikhil",
				Password: "123",
				Type:     "regular",
			},
			"Rewak": {
				Username: "Rewak",
				Password: "222",
				Type:     "admin",
			},
			"Prince": {
				Username: "Prince",
				Password: "333",
				Type:     "regular",
			},
			"Abhishek": {
				Username: "Abhishek",
				Password: "111",
				Type:     "admin",
			},
		},
	}
}

func (app *Application) generateJWT(u models.User) (string, error) {
	var mySigningKey = []byte(app.SecretKey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["username"] = u.Username
	claims["type"] = u.Type

	return token.SignedString(mySigningKey)
}

func readCSV(filepath string) ([][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	return reader.ReadAll()
}

// comment for reviewers: this func will work if last line is empty in the csv file
// otherwise an unnecessary empty line was being created when I try implement it without an empty line in the end
func writeBookToCSV(filepath string, data []models.Book) error {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	startPos, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	for _, book := range data {
		record := []string{book.Name, book.Author, strconv.Itoa(book.Year)}
		err := writer.Write(record)
		if err != nil {
			return err
		}
	}

	if err := file.Truncate(startPos); err != nil {
		return err
	}

	return nil
}

func deleteBookFromCSV(filepath string, bookName string) error {
	file, err := os.OpenFile(filepath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return err
	}

	var filteredRecords [][]string
	for _, record := range records {
		if len(record) > 0 && !strings.EqualFold(record[0], bookName) {
			filteredRecords = append(filteredRecords, record)
		}
	}

	if len(filteredRecords) == len(records) {
		return models.ErrBookNotExists
	}

	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range filteredRecords {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func internalServerError(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
}

func unauthorizedError(c echo.Context) error {
	return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user unauthorized"})
}

func jsonError(err error) map[string]string {
	return map[string]string{"error": err.Error()}
}
