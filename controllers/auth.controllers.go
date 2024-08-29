package controllers

import (
	"ecommerce/configs"
	"ecommerce/helpers"
	"ecommerce/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Registration struct {
	Email        string `json:"email" validate:"required,email"`
	Name         string `json:"name" validate:"required"`
	Password     string `json:"password" validate:"required"`
	PasswordConf string `json:"confirmPassword" validate:"required"`
	Language     string `json:"language"`
}
type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserUpdatePayload struct {
	Name            string `json:"name"`
	DefaultLanguage string `json:"defaultLanguage"`
	Photo           string `json:"photo"`
}

func Register(c *fiber.Ctx) error {

	var payload *Registration

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error(), "success": false})
	}

	if errors := helpers.ValidateStruct(payload); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": errors, "success": false})

	}

	if payload.Password != payload.PasswordConf {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Passwords do not match", "success": false})
	}

	var user models.User

	db := configs.DB

	alreadyRegistered := db.Where("email = ?", strings.ToLower(payload.Email)).First(&user)

	if alreadyRegistered.RowsAffected > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "User with that email already exists", "success": false})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": err.Error(), "success": false})
	}

	newUser := &models.User{
		Name:            payload.Name,
		Email:           strings.ToLower(payload.Email),
		Password:        string(hashedPassword),
		DefaultLanguage: payload.Language,
	}

	result := db.Create(&newUser)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Something bad happened"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User successfully registered", "data": fiber.Map{"id": newUser.ID}, "success": true})

}
func Login(c *fiber.Ctx) error {

	var payload *LoginPayload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error(), "success": false})
	}

	if errors := helpers.ValidateStruct(payload); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errors, "success": false})

	}

	db := configs.DB
	userExist := models.User{}
	result := db.First(&userExist, "email = ?", strings.ToLower(payload.Email))

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Email or Password is incorrect", "success": false})
	}

	if userExist.IsBlocked {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Account is blocked! Please contact support", "success": false})
	}

	comparisonError := bcrypt.CompareHashAndPassword([]byte(userExist.Password), []byte(payload.Password))

	if comparisonError != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Email or Password is incorrect", "success": false})
	}

	refreshToken, err := helpers.GenerateToken(helpers.TokenClaims{
		UserId: userExist.ID.String(),
		Role:   userExist.Role,
		Exp:    time.Now().Add(time.Hour * 168).Unix(), // 7 days
	})

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Couldn't login user", "success": false})
	}

	cookie := fiber.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 168), // Set expiration
		HTTPOnly: true,                            // JavaScript can't access the cookie
		Secure:   true,                            // Only send over HTTPS
		SameSite: fiber.CookieSameSiteStrictMode,  // Protect against CSRF
	}

	c.Cookie(&cookie)

	userExist.IsLoggedIn = true
	db.Save(&userExist)

	accessToken, err := helpers.GenerateToken(helpers.TokenClaims{
		UserId: userExist.ID.String(),
		Role:   userExist.Role,
		Exp:    time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Couldn't login user", "success": false})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User successfully logged in", "data": fiber.Map{"accessToken": accessToken, "user": fiber.Map{"id": userExist.ID, "name": userExist.Name, "email": userExist.Email, "role": userExist.Role}}, "success": true})

}
func GenerateToken(c *fiber.Ctx) error {

	cookies := c.Cookies("refreshToken")

	if cookies == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Please login", "success": false})
	}

	validToken, err := helpers.ParseToken(cookies)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Please login", "success": false})
	}

	userId := validToken.UserId
	parsedUuid, err := uuid.Parse(userId)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Please login", "success": false})
	}

	db := configs.DB
	user := models.User{}
	result := db.Select("is_blocked", "id", "role", "valid_logged_in_at").First(&user, "id = ?", parsedUuid)

	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Please login", "success": false})
	}

	if user.IsBlocked {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Please login", "success": false})
	}

	if user.ValidLoggedInAt != nil && validToken.Exp < user.ValidLoggedInAt.Unix() {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Please login", "success": false})
	}

	accessToken, err := helpers.GenerateToken(helpers.TokenClaims{
		UserId: user.ID.String(),
		Role:   user.Role,
		Exp:    time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Please login", "success": false})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Token successfully generated", "data": fiber.Map{"accessToken": accessToken}, "success": true})

}
func LogoutUser(c *fiber.Ctx) error {

	userId := c.Locals("userId")
	parsedUuid, err := uuid.Parse(userId.(string))

	logoutEveryDevice := c.Query("logoutEveryDevice", "false")

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Please login", "success": false})
	}

	db := configs.DB
	user := models.User{}
	result := db.Select("id").First(&user, "id = ?", parsedUuid)

	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Please login", "success": false})
	}

	updateUsers := models.User{
		IsLoggedIn: false,
	}

	if logoutEveryDevice == "true" {
		now := time.Now()
		updateUsers.ValidLoggedInAt = &now
	}

	if err := db.Model(&user).Updates(updateUsers).Error; err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Couldn't logged out user", "success": false})
	}

	cookie := fiber.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Expires:  time.Now().Add(time.Hour * -1),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
	}

	c.Cookie(&cookie)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User successfully logged out", "success": true})

}
func GetProfile(c *fiber.Ctx) error {

	userId := c.Locals("userId")
	parsedUuid, err := uuid.Parse(userId.(string))

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Please login", "success": false})
	}

	db := configs.DB
	user := models.User{}
	result := db.Select("id", "name", "email", "role", "IsBlocked").First(&user, "id = ?", parsedUuid)

	if result.Error != nil {
		c.ClearCookie("refreshToken")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Please login", "success": false})
	}

	if user.IsBlocked {
		c.ClearCookie("refreshToken")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Account blocked! Please contact support", "success": false})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User profile fetched successfully", "data": fiber.Map{
		"name":           user.Name,
		"email":          user.Email,
		"role":           user.Role,
		"isBlocked":      user.IsBlocked,
		"defaultLang":    user.DefaultLanguage,
		"photo":          user.Photo,
		"createdAt":      user.CreatedAt,
		"defaultAddress": user.DefaultAddress,
		"id":             user.ID,
	}, "success": true})

}
func UpdateProfile(c *fiber.Ctx) error {

	var payload *UserUpdatePayload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error(), "success": false})
	}

	if errors := helpers.ValidateStruct(payload); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errors, "success": false})

	}

	userId := c.Locals("userId")
	parsedUuid, err := uuid.Parse(userId.(string))

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Please login", "success": false})
	}

	db := configs.DB
	user := models.User{}
	result := db.First(&user, "id = ?", parsedUuid)

	if result.Error != nil {
		c.ClearCookie("refreshToken")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Please login", "success": false})
	}

	userUpdate := models.User{
		Name:            payload.Name,
		DefaultLanguage: payload.DefaultLanguage,
		Photo:           payload.Photo,
	}

	if err := db.Model(&user).Updates(userUpdate).Error; err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Couldn't update user profile", "success": false})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User profile updated successfully", "success": true})

}
