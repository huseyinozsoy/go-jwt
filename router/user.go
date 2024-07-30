package router

import (
	"math/rand"
	"os"

	db "github.com/huseyinozsoy/go-jwt/database"

	"github.com/gofiber/fiber/v2"

	"github.com/huseyinozsoy/go-jwt/models"
	"github.com/huseyinozsoy/go-jwt/utils"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("PRIV_KEY"))

// SetupUserRoutes func sets up all the user routes
func SetupUserRoutes() {
	USER.Post("/signup", CreateUser) // Sign Up a user

	USER.Post("/signin", LoginUser) // Sign In a user

	// privUser handles all the private user routes that requires authentication
	privUser := USER.Group("/private")
	privUser.Use(utils.SecureAuth()) // middleware to secure all routes for this group
	privUser.Get("/user", GetUserData)
}

// CreateUser route registers a User into the database
func CreateUser(c *fiber.Ctx) error {
	u := new(models.User)

	if err := c.BodyParser(u); err != nil {
		return c.JSON(fiber.Map{
			"error": true,
			"input": "Please review your input",
		})
	}

	// validate if the email, username and password are in correct format
	errors := utils.ValidateRegister(u)
	if errors.Status {
		return c.JSON(errors)
	}

	if count := db.DB.Where(&models.User{Email: u.Email}).First(new(models.User)).RowsAffected; count > 0 {
		errors.Status, errors.Error = true, "Email is already registered"
	}
	if count := db.DB.Where(&models.User{Username: u.Username}).First(new(models.User)).RowsAffected; count > 0 {
		errors.Status, errors.Error = true, "Username is already registered"
	}
	if errors.Status {
		return c.JSON(errors)
	}

	// Hashing the password with a random salt
	password := []byte(u.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(
		password,
		rand.Intn(bcrypt.MaxCost-bcrypt.MinCost)+bcrypt.MinCost,
	)

	if err != nil {
		panic(err)
	}
	u.Password = string(hashedPassword)
	u.BeforeCreate(db.DB)

	if err := db.DB.Create(&u).Error; err != nil {
		return c.JSON(fiber.Map{
			"error":   true,
			"general": "Something went wrong, please try again later.",
		})
	}

	// setting up the authorization cookies
	accessToken := utils.GenerateTokens(u.ID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": accessToken,
	})
}

// LoginUser route logins a user in the app
func LoginUser(c *fiber.Ctx) error {
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	input := new(LoginInput)

	if err := c.BodyParser(input); err != nil {
		return c.JSON(fiber.Map{"error": true, "input": "Please review your input"})
	}

	// check if a user exists
	u := new(models.User)
	if res := db.DB.Where(
		&models.User{Email: input.Identity}).Or(
		&models.User{Username: input.Identity},
	).First(&u); res.RowsAffected <= 0 {
		return c.JSON(fiber.Map{"error": true, "general": "Invalid Credentials."})
	}

	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); err != nil {
		return c.JSON(fiber.Map{"error": true, "general": "Invalid Credentials."})
	}

	// setting up the authorization cookies
	accessToken := utils.GenerateTokens(u.ID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": accessToken,
	})
}

// GetUserData returns the details of the user signed in
func GetUserData(c *fiber.Ctx) error {
	id := c.Locals("id")

	u := new(models.User)
	if res := db.DB.Where("id = ?", id).First(&u); res.RowsAffected <= 0 {
		return c.JSON(fiber.Map{"error": true, "general": "Cannot find the User"})
	}

	return c.JSON(u)
}
