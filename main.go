package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var db *sqlx.DB

func main() {
	var err error

	db, err = sqlx.Open("mysql", "root:IntelliP24@tcp(localhost:3306)/namedeelukrak_db")
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	app.Post("/signup", func(c *fiber.Ctx) error {
		request := SignupRequest{}
		err := c.BodyParser(&request)
		if err != nil {
			return err
		}

		if request.Email == "" || request.Username == "" || request.Password == "" {
			return fiber.ErrUnprocessableEntity
		}

		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
		if err != nil {
			return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		}

		query := "INSERT member (email, username, password) values (?, ?, ?)"
		result, err := db.Exec(query, request.Email, request.Username, string(password))
		if err != nil {
			return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		}
		id, err := result.LastInsertId()
		if err != nil {
			return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		}

		member := Member{
			Id:       int(id),
			Email:    request.Email,
			Password: request.Password,
		}

		return c.Status(fiber.StatusCreated).JSON(member)
	})

	app.Static("/", "./wwwroot")
	app.Listen(":4542")

}

type Member struct {
	Id       int    `db:"id" json:"id"`
	Email    string `db:"email json:"email"`
	Username string `db:"username json:"username"`
	Password string `db:"password json:"password"`
}

type SignupRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
