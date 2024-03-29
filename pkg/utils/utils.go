package utils

import (
	"regexp"
)

func ExtractEmails(message string) ([]string, error) {
	pattern := `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return re.FindAllString(message, -1), nil
}

func IsValidEmail(email string) (bool, error) {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	emailRegex, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}
	return emailRegex.MatchString(email), nil
}

func AreValidEmails(emails []string) (bool, error) {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	emailRegex, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}
	for _, email := range emails {
		if !emailRegex.MatchString(email) {
			return false, nil
		}
	}
	return true, nil
}
