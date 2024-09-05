package helpers

import "fmt"

func GenerateOtp() string {
	return "123456"
}

func SendOTP(otp string, mobile string) error {
	fmt.Println("OTP: ", otp, "Mobile: ", mobile)
	return nil
}
