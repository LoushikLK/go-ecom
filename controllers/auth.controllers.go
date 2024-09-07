package controllers

import (
	"ecommerce/configs"
	"ecommerce/helpers"
	"ecommerce/models"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AccountRegistrationPayload struct {
	Mobile      string `json:"mobile" validate:"required,number,max=10,min=10"`
	CountryCode string `json:"country_code" validate:"required,max=3"`
	Fcm         string `json:"fcm" validate:"required"`
	Platform    string `json:"platform" validate:"required"`
	Language    string `json:"language"`
}
type ResendVerificationOTPPayload struct {
	Mobile string `json:"mobile" validate:"required,number,max=10,min=10"`
}
type VerifyAccountPayload struct {
	Mobile string `json:"mobile" validate:"required,number,max=10,min=10"`
	Otp    string `json:"otp" validate:"required"`
	Fcm    string `json:"fcm" validate:"required"`
}
type LoginPayload struct {
	Mobile string `json:"mobile" validate:"required,number,max=10,min=10"`
}
type LoginVerifyPayload struct {
	Mobile   string `json:"mobile" validate:"required,number,max=10,min=10"`
	Otp      string `json:"otp" validate:"required"`
	FCM      string `json:"fcm" validate:"required"`
	Platform string `json:"platform" validate:"required"`
}
type GenerateAccessTokenPayload struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func Register(c *fiber.Ctx) error {

	var payload *AccountRegistrationPayload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error(), "success": false})
	}

	if err := helpers.ValidateStruct(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err, "success": false})
	}

	db := configs.DB

	user := models.Account{}

	// Check if user already exists
	alreadyRegistered := db.First(&user, "mobile = ?", payload.Mobile)

	if alreadyRegistered.RowsAffected > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "User with that phone number already exists! Please login!", "success": false})
	}

	newUser := models.Account{
		Mobile:      payload.Mobile,
		CountryCode: payload.CountryCode,
	}

	// Create user
	result := db.Create(&newUser)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "User with that phone number or email already exists! Please login!", "success": false})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Something bad happened", "success": false})
	}

	//save user login
	userLogin := models.UserLogin{
		FCM:       payload.Fcm,
		Platform:  payload.Platform,
		AccountID: newUser.ID,
		Lang:      payload.Language,
	}

	userLoginResult := db.Create(&userLogin)

	if userLoginResult.RowsAffected == 0 {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Could not create user login", "success": false})
	}

	if userLoginResult.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Something bad happened", "success": false})
	}

	// Generate OTP
	otp := helpers.GenerateOtp()

	// This is the first time it will be created after registration after we only need to update the OTP table
	createOtp := db.Create(&models.UserOtp{
		AccountID:       newUser.ID,
		Otp:             otp,
		ExpiredDateTime: time.Now().Add(10 * time.Minute), // 10 minutes
	})

	if createOtp.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Something bad happened", "success": false})
	}

	err := helpers.SendOTP(otp, newUser.Mobile)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to send OTP", "success": false})
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User successfully registered", "data": fiber.Map{"id": newUser.ID}, "success": true})
}
func ResendVerificationOtp(c *fiber.Ctx) error {

	var payload *ResendVerificationOTPPayload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error(), "success": false})
	}

	if errors := helpers.ValidateStruct(payload); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": errors, "success": false})

	}

	db := configs.DB

	user := models.Account{}

	alreadyRegistered := db.First(&user, "mobile = ?", payload.Mobile)

	if alreadyRegistered.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Something bad happened", "success": false})
	}

	if alreadyRegistered.RowsAffected == 0 {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "User is not registered. Please register!", "success": false})
	}

	if user.IsMobileVerified {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "User is already verified. Please login!", "success": false})
	}

	if user.IsBlacklisted {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "User is blacklisted. Please contact support!", "success": false})
	}

	if user.IsBlocked {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "User is blocked. Please contact support!", "success": false})
	}

	// Generate OTP
	otp := helpers.GenerateOtp()

	userOtp := models.UserOtp{}

	createOtp := db.FirstOrCreate(&userOtp, models.UserOtp{AccountID: user.ID})

	if createOtp.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Something bad happened", "success": false})
	}

	userOtp.Otp = otp
	userOtp.ExpiredDateTime = time.Now().Add(5 * time.Minute)

	result := db.Save(&userOtp)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Something bad happened", "success": false})
	}

	err := helpers.SendOTP(otp, user.Mobile)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to send OTP", "success": false})
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "OTP successfully sent", "data": fiber.Map{"id": user.ID}, "success": true})
}
func VerifyAccountRegistration(c *fiber.Ctx) error {

	var payload *VerifyAccountPayload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error(), "success": false})
	}

	if errors := helpers.ValidateStruct(payload); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": errors, "success": false})

	}

	db := configs.DB

	user := models.Account{}

	alreadyRegistered := db.First(&user, "mobile = ?", payload.Mobile)

	if alreadyRegistered.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User is not registered. Please register!", "success": false})
	}

	if alreadyRegistered.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Something bad happened", "success": false})
	}

	if user.IsMobileVerified {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "User is already verified. Please login!", "success": false})
	}

	if user.IsBlacklisted {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "User is blacklisted. Please contact support!", "success": false})
	}

	if user.IsBlocked {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "User is blocked. Please contact support!", "success": false})
	}

	userLoginModel := models.UserLogin{}

	loginResult := db.First(&userLoginModel, "account_id = ? AND fcm = ?", user.ID, payload.Fcm)

	if loginResult.RowsAffected == 0 {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "You are tried to login from different device than registered account. Please use registered device!", "success": false})
	}

	otpModel := models.UserOtp{}

	otpResult := db.First(&otpModel, "account_id = ?", user.ID)

	if otpResult.RowsAffected == 0 {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "No OTP found", "success": false})
	}

	if otpResult.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": otpResult.Error.Error(), "success": false})
	}

	if otpModel.IsExpired {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "OTP Expired!", "success": false})
	}

	if otpModel.ExpiredDateTime.UnixMilli() < time.Now().UnixMilli() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "OTP Expired!", "success": false})
	}

	if otpModel.Otp != payload.Otp {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid OTP!", "success": false})
	}

	userUpdate := models.Account{}

	result := db.Model(&userUpdate).Where("id = ?", user.ID).Updates(&models.Account{IsMobileVerified: true})

	if result.Error != nil {
		log.Printf("Error updating user: %v", result.Error)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Something bad happened", "success": false})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User verified successfully", "data": fiber.Map{"id": user.ID}, "success": true})
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
	userExist := models.Account{}
	result := db.First(&userExist, "mobile = ?", strings.ToLower(payload.Mobile))

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Something bad happened", "success": false})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Account does not exist. Please register!", "success": false})
	}

	if userExist.IsBlocked {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": userExist.IsBlockedReason, "success": false})
	}
	if userExist.IsBlacklisted {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": userExist.IsBlacklistedReason, "success": false})
	}
	if !userExist.IsMobileVerified {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "You are not verified! Please verify your account!", "data": fiber.Map{"id": userExist.ID, "isMobileVerified": userExist.IsMobileVerified}, "success": false})
	}

	// Generate OTP
	otp := helpers.GenerateOtp()

	userOtp := models.UserOtp{}

	createOtp := db.FirstOrCreate(&userOtp, models.UserOtp{AccountID: userExist.ID})

	if createOtp.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Something bad happened", "success": false})
	}

	userOtp.Otp = otp
	userOtp.ExpiredDateTime = time.Now().Add(5 * time.Minute)

	otpSaveResult := createOtp.Save(&userOtp)

	if otpSaveResult.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Something bad happened", "success": false})
	}

	err := helpers.SendOTP(otp, userExist.Mobile)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to send OTP", "success": false})
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "OTP successfully sent", "success": true})

}
func LoginVerifyOTP(c *fiber.Ctx) error {

	var payload *LoginVerifyPayload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error(), "success": false})
	}

	if errors := helpers.ValidateStruct(payload); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errors, "success": false})
	}

	db := configs.DB
	userExist := models.Account{}
	result := db.First(&userExist, "mobile = ?", strings.ToLower(payload.Mobile))

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Something bad happened", "success": false})
	}

	if userExist.IsBlocked {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": userExist.IsBlockedReason, "success": false})
	}
	if userExist.IsBlacklisted {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": userExist.IsBlacklistedReason, "success": false})
	}
	if !userExist.IsMobileVerified {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "You are not verified! Please verify your account!", "success": false})
	}

	//find user otp and match the otp
	userOtp := models.UserOtp{}
	result = db.First(&userOtp, "account_id = ? AND otp = ?", userExist.ID, payload.Otp)
	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Something bad happened", "success": false})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Invalid OTP. Please try again!", "success": false})
	}

	if userOtp.IsExpired {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "OTP has expired. Please try again!", "success": false})
	}

	if time.Now().After(userOtp.ExpiredDateTime) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "OTP has expired. Please try again!", "success": false})
	}

	userOtp.IsExpired = true
	db.Save(&userOtp)

	userLogin := models.UserLogin{}

	refreshToken, err := helpers.GenerateToken(helpers.TokenClaims{
		UserId: userExist.ID,
		Exp:    time.Now().Add(time.Hour * 168).Unix(), // 7 days
	})

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Couldn't login user", "success": false})
	}

	userLoginUpdate := db.FirstOrCreate(&userLogin, models.UserLogin{AccountID: userExist.ID, FCM: payload.FCM, Platform: payload.Platform})

	if userLoginUpdate.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Couldn't login user", "success": false})
	}

	userExist.IsLoggedIn = true
	db.Save(&userExist)

	userLogin.RefreshToken = refreshToken
	userLogin.IsActive = true
	db.Save(&userLogin)

	//update other login of the same platform with is_active = false
	db.Model(&models.UserLogin{}).Where("id != ? AND account_id = ? AND platform = ? AND is_active = ?", userLogin.ID, userExist.ID, payload.Platform, true).Update("is_active", false)

	accessToken, err := helpers.GenerateToken(helpers.TokenClaims{
		UserId: userExist.ID,
		Exp:    time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})

	if err != nil {
		log.Printf("Error generating access token: %v", err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Couldn't login user", "success": false})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User successfully logged in", "data": fiber.Map{"accessToken": accessToken, "refreshToken": refreshToken, "user": fiber.Map{"id": userExist.ID, "name": userExist.Name, "email": userExist.Email}}, "success": true})

}
func GenerateToken(c *fiber.Ctx) error {

	var payload *GenerateAccessTokenPayload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error(), "success": false})
	}

	if errors := helpers.ValidateStruct(payload); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errors, "success": false})
	}

	validToken, err := helpers.ParseToken(payload.RefreshToken)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Please login", "success": false})
	}

	userId := validToken.UserId

	db := configs.DB
	user := models.Account{}
	result := db.Select("is_blocked", "id", "role", "valid_logged_in_at").First(&user, "id = ?", userId)

	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Please login", "success": false})
	}

	if user.IsBlocked {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Please login", "success": false})
	}

	if !user.IsBlacklisted {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Please login", "success": false})
	}

	accessToken, err := helpers.GenerateToken(helpers.TokenClaims{
		UserId: user.ID,
		Exp:    time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Please login", "success": false})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Token successfully generated", "data": fiber.Map{"accessToken": accessToken}, "success": true})

}
func LogoutUser(c *fiber.Ctx) error {

	userId := c.Locals("userId")

	db := configs.DB
	user := models.Account{}
	result := db.Select("id").First(&user, "id = ?", userId)

	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Please login", "success": false})
	}

	updateUsers := models.Account{
		IsLoggedIn: false,
	}

	if err := db.Model(&user).Updates(updateUsers).Error; err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Couldn't logged out user", "success": false})
	}

	db.Model(&models.UserLogin{}).Where("account_id = ?", user.ID).Update("is_active", false)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User successfully logged out", "success": true})

}
