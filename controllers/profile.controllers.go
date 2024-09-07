package controllers

import (
	"ecommerce/configs"
	"ecommerce/helpers"
	"ecommerce/models"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UpdateAccountPayload struct {
	Name  string  `json:"name"`
	Lang  string  `json:"lang"`
	Photo string  `json:"photo"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}

type UpdateEmailPayload struct {
	Email string `json:"email" validate:"required,email"`
}
type VerifyEmailPayload struct {
	Token string `json:"token" validate:"required"`
}

type AddAddressPayload struct {
	IsDefault         bool    `json:"is_default" validate:"boolean"`
	FullName          string  `json:"full_name" validate:"omitempty,min=3,max=100"`
	PhoneNumber       string  `json:"phone_number" validate:"omitempty,number,max=10,min=10"`
	CountryCode       string  `json:"country_code" validate:"omitempty,number,max=3"`
	AddressLine1      string  `json:"address_line1" validate:"omitempty,min=15,max=255"`
	AddressLine2      string  `json:"address_line2" validate:"omitempty,min=15,max=255"`
	PostalCode        string  `json:"postal_code" validate:"required,number,max=6,min=6"`
	City              string  `json:"city" validate:"required,min=3,max=100"`
	State             string  `json:"state" validate:"required,min=3,max=100"`
	Country           string  `json:"country " validate:"omitempty,min=3,max=100"`
	IsShippingAddress bool    `json:"is_shipping_address" validate:"omitempty,boolean"`
	IsBillingAddress  bool    `json:"is_billing_address" validate:"omitempty,boolean"`
	AddressTitle      string  `json:"address_title" validate:"omitempty,min=3,max=100"`
	Lat               float64 `json:"lat" validate:"omitempty,number"`
	Long              float64 `json:"long" validate:"omitempty,number"`
}
type UpdateAddressPayload struct {
	IsDefault         bool    `json:"is_default" validate:"omitempty,boolean"`
	FullName          string  `json:"full_name" validate:"omitempty,min=3,max=100"`
	PhoneNumber       string  `json:"phone_number" validate:"omitempty,number,max=10,min=10"`
	CountryCode       string  `json:"country_code" validate:"omitempty,number,max=3"`
	AddressLine1      string  `json:"address_line1" validate:"omitempty,min=15,max=255"`
	AddressLine2      string  `json:"address_line2" validate:"omitempty,min=15,max=255"`
	PostalCode        string  `json:"postal_code" validate:"omitempty,number,max=6,min=6"`
	City              string  `json:"city" validate:"omitempty,min=3,max=100"`
	State             string  `json:"state" validate:"omitempty,min=3,max=100"`
	Country           string  `json:"country" validate:"omitempty,min=3,max=100"`
	IsShippingAddress bool    `json:"is_shipping_address" validate:"omitempty,boolean"`
	IsBillingAddress  bool    `json:"is_billing_address" validate:"omitempty,boolean"`
	AddressTitle      string  `json:"address_title" validate:"omitempty,min=3,max=100"`
	Lat               float64 `json:"lat" validate:"omitempty,number"`
	Long              float64 `json:"long" validate:"omitempty,number"`
}

type GetAddressQuery struct {
	Limit int    `json:"limit" validate:"omitempty,number,min=1,max=100"`
	Page  int    `json:"page" validate:"omitempty,number"`
	Sort  string `json:"sort" validate:"omitempty,oneof=asc desc"`
}

func GetProfile(c *fiber.Ctx) error {

	userId := c.Locals("userId")

	db := configs.DB
	user := models.Account{}
	result := db.Select("id", "name", "email", "is_blocked", "mobile", "is_blacklisted", "lang", "country_code", "is_mobile_verified", "is_email_verified", "profile_image").First(&user, "id = ?", userId)

	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Please login", "success": false})
	}

	if user.IsBlocked {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Account blocked! Please contact support", "success": false})
	}

	if user.IsBlacklisted {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Account blacklisted! Please contact support", "success": false})
	}

	if !user.IsMobileVerified {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unverified account please verify", "success": false})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User profile fetched successfully", "data": fiber.Map{
		"name":               user.Name,
		"id":                 user.ID,
		"email":              user.Email,
		"mobile":             user.Mobile,
		"is_blocked":         user.IsBlocked,
		"is_blacklisted":     user.IsBlacklisted,
		"is_mobile_verified": user.IsMobileVerified,
		"is_email_verified":  user.IsEmailVerified,
		"lang":               user.Lang,
		"country_code":       user.CountryCode,
		"profile_image":      user.ProfileImage,
	}, "success": true})

}
func UpdateProfile(c *fiber.Ctx) error {

	var payload *UpdateAccountPayload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error(), "success": false})
	}

	if errors := helpers.ValidateStruct(payload); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errors, "success": false})

	}

	userId := c.Locals("userId")

	db := configs.DB
	user := models.Account{}
	result := db.First(&user, "id = ?", userId)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Something went wrong", "success": false})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized access. Please login", "success": false})
	}

	userUpdate := models.Account{}

	if payload.Name != "" {
		userUpdate.Name = payload.Name
	}
	if payload.Photo != "" {
		userUpdate.ProfileImage = payload.Photo
	}

	if payload.Lang != "" {
		userUpdate.Lang = payload.Lang
	}

	if payload.Lat != 0 {
		userUpdate.Lat = payload.Lat
	}

	if payload.Long != 0 {
		userUpdate.Long = payload.Long
	}

	if err := db.Model(&user).Updates(&userUpdate).Error; err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Couldn't update user profile", "success": false})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User profile updated successfully", "success": true})

}

func UpdateEmail(c *fiber.Ctx) error {

	var payload *UpdateEmailPayload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error(), "success": false})
	}

	if errors := helpers.ValidateStruct(payload); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errors, "success": false})
	}

	userId := c.Locals("userId")

	db := configs.DB
	userExist := models.Account{}
	result := db.First(&userExist, "id = ?", userId)

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

	// Generate token based on email and send email
	token, err := helpers.GenerateToken(helpers.TokenClaims{
		UserId: userExist.ID,
		Email:  payload.Email,
		Exp:    time.Now().Add(5 * time.Minute).Unix(), // 5 minutes from now
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Something bad happened on server", "success": false})
	}

	// Send email
	err = helpers.SendEmail(token, userExist.Email)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to send email", "success": false})
	}

	userExist.Email = payload.Email
	userExist.IsEmailVerified = false

	db.Save(&userExist)

	// Return response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Email updated and verification email sent", "success": true})

}
func VerifyEmail(c *fiber.Ctx) error {

	var payload *VerifyEmailPayload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error(), "success": false})
	}

	if errors := helpers.ValidateStruct(payload); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errors, "success": false})
	}

	parseToken, err := helpers.ParseToken(payload.Token)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error(), "success": false})
	}

	userId := parseToken.UserId

	db := configs.DB
	userExist := models.Account{}
	result := db.First(&userExist, "id = ?", userId)

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

	if userExist.Email != parseToken.Email {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token. Failed to verify email!", "success": false})
	}

	userExist.IsEmailVerified = true

	db.Save(&userExist)

	// Return response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Email verified successfully", "success": true})

}

func AddAddress(c *fiber.Ctx) error {

	var payload *AddAddressPayload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error(), "success": false})
	}

	if errors := helpers.ValidateStruct(payload); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errors, "success": false})
	}

	userId := c.Locals("userId")

	db := configs.DB
	user := models.Account{}
	result := db.First(&user, "id = ?", userId)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Account does not exist. Please register!", "success": false})
	}

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Something bad happened", "success": false})
	}

	userAddress := models.Address{
		AccountID: user.ID,
	}

	if payload.FullName != "" {
		userAddress.FullName = payload.FullName
	} else {
		userAddress.FullName = user.Name
	}

	if payload.PhoneNumber != "" {
		userAddress.PhoneNumber = payload.PhoneNumber
	} else {
		userAddress.PhoneNumber = user.Mobile
	}

	if payload.CountryCode != "" {
		userAddress.CountryCode = payload.CountryCode
	} else {
		userAddress.CountryCode = user.CountryCode
	}

	if payload.AddressLine1 != "" {
		userAddress.AddressLine1 = payload.AddressLine1
	}

	if payload.AddressLine2 != "" {
		userAddress.AddressLine2 = payload.AddressLine2
	}

	userAddress.PostalCode = payload.PostalCode
	userAddress.City = payload.City
	userAddress.State = payload.State

	if payload.Country != "" {
		userAddress.Country = payload.Country
	} else {
		userAddress.Country = "India"
	}

	if payload.IsShippingAddress {
		userAddress.IsShippingAddress = true
	}

	if payload.IsBillingAddress {
		userAddress.IsBillingAddress = true
	}

	if payload.AddressTitle != "" {
		userAddress.AddressTitle = payload.AddressTitle
	}

	if payload.Lat != 0 {
		userAddress.Lat = payload.Lat
	}

	if payload.Long != 0 {
		userAddress.Long = payload.Long
	}

	transaction := db.Begin()

	if payload.IsDefault {
		userAddress.IsDefault = true
		result := transaction.Model(&userAddress).Unscoped().Where("account_id = ?", user.ID).Update("is_default", false)

		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			transaction.Rollback()
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Adding address failed. Please try again!", "success": false})
		}
	}

	if err := transaction.Create(&userAddress).Error; err != nil {
		transaction.Rollback()
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Adding address failed. Please try again!", "success": false})
	}

	if err := transaction.Commit().Error; err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Adding address failed. Please try again!", "success": false})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Address added successfully", "success": true})

}

func UpdateAddress(c *fiber.Ctx) error {

	var payload *UpdateAddressPayload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error(), "success": false})
	}

	if errors := helpers.ValidateStruct(payload); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errors, "success": false})

	}

	userId := c.Locals("userId")
	addressId := c.Params("addressId")

	if addressId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid address id", "success": false})
	}

	db := configs.DB
	address := models.Address{}
	result := db.First(&address, "id = ?  AND account_id = ? AND is_deleted = false", addressId, userId)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Something bad happened while fetching address", "success": false})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Address does not exist. Please add address!", "success": false})
	}

	if payload.FullName != "" {
		address.FullName = payload.FullName
	}

	if payload.PhoneNumber != "" {
		address.PhoneNumber = payload.PhoneNumber
	}

	if payload.CountryCode != "" {
		address.CountryCode = payload.CountryCode
	}

	if payload.AddressLine1 != "" {
		address.AddressLine1 = payload.AddressLine1
	}

	if payload.AddressLine2 != "" {
		address.AddressLine2 = payload.AddressLine2
	}

	if payload.PostalCode != "" {
		address.PostalCode = payload.PostalCode
	}

	if payload.City != "" {
		address.City = payload.City
	}

	if payload.State != "" {
		address.State = payload.State
	}

	if payload.Country != "" {
		address.Country = payload.Country
	}

	if payload.IsShippingAddress {
		address.IsShippingAddress = true
	}

	if payload.IsBillingAddress {
		address.IsBillingAddress = true
	}

	if payload.AddressTitle != "" {
		address.AddressTitle = payload.AddressTitle
	}

	if payload.Lat != 0 {
		address.Lat = payload.Lat
	}

	if payload.Long != 0 {
		address.Long = payload.Long
	}

	transaction := db.Begin()

	if payload.IsDefault {
		address.IsDefault = true
		result := transaction.Model(&address).Unscoped().Where("account_id = ?", userId).Update("is_default", false)
		if result.Error != nil {
			transaction.Rollback()
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Update address failed. Please try again!", "success": false})
		}
	}

	if err := transaction.Updates(&address).Error; err != nil {
		transaction.Rollback()
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Update address failed. Please try again!", "success": false})
	}

	if err := transaction.Commit().Error; err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Update address failed. Please try again!", "success": false})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Address updated successfully", "success": true})
}
func DeleteAddress(c *fiber.Ctx) error {

	userId := c.Locals("userId")
	addressId := c.Params("addressId")

	if addressId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid address id", "success": false})
	}

	db := configs.DB
	address := models.Address{}
	result := db.First(&address, "id = ?  AND account_id = ? AND is_deleted = false", addressId, userId)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Something bad happened while fetching address", "success": false})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Address does not exist. Please try again!", "success": false})
	}

	if address.IsDefault {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "You cannot delete default address!", "success": false})
	}

	if err := db.Model(&address).Updates(&models.Address{IsDeleted: true, DeletedAt: time.Now()}).Error; err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Delete address failed. Please try again!", "success": false})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Address deleted successfully", "success": true})
}
func GetAllAddresses(c *fiber.Ctx) error {

	var payload GetAddressQuery

	if err := c.QueryParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error(), "success": false})
	}

	if errors := helpers.ValidateStruct(payload); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errors, "success": false})
	}

	userId := c.Locals("userId")

	db := configs.DB
	var addresses []models.Address

	defaultQuery := db.Model(&addresses).Where("account_id = ? AND is_deleted = false", userId)

	var Limit int
	var Page int
	var Sort string

	if payload.Limit == 0 {
		Limit = 10
	} else if payload.Limit > 100 {
		Limit = 100
	} else if payload.Limit < 1 {
		Limit = 1
	} else {
		Limit = payload.Limit
	}

	if payload.Page <= 0 {
		Page = 1
	}

	if payload.Sort == "asc" {
		Sort = "created_at asc"
	} else {
		Sort = "created_at desc"
	}

	if Limit != 0 {
		defaultQuery = defaultQuery.Limit(Limit)
	}

	if Page != 0 {
		defaultQuery = defaultQuery.Offset((Page - 1) * Limit)
	}

	if Sort != "" {
		defaultQuery = defaultQuery.Order(Sort)
	}

	selectQuery := []string{"id", "is_default", "full_name", "phone_number", "country_code", "address_line1", "address_line2", "city", "state", "country", "postal_code", "is_shipping_address", "is_billing_address", "address_title", "lat", "long", "created_at", "updated_at"}

	result := defaultQuery.Select(selectQuery).Find(&addresses)

	if result.Error != nil {
		log.Println(result.Error)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Something bad happened while fetching address", "success": false})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Address fetched successfully", "data": addresses, "count": len(addresses), "total": result.RowsAffected, "page": payload.Page, "limit": payload.Limit, "success": true})
}
