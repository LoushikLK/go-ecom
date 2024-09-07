package helpers

import "fmt"

func SendEmail(email string, text string) error {
	fmt.Println("Email: ", email, "Text: ", text)
	return nil
}
