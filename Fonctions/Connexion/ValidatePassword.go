package Connexion

import "unicode"

func ValidatePassword(password string) bool {
	hasUpper := false
	hasSpecial := false
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecial = true
		}
	}
	return len(password) >= 8 && hasUpper && hasSpecial
}
