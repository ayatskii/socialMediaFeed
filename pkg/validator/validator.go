package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) IsValid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(field, message string) {
	if _, exists := v.Errors[field]; !exists {
		v.Errors[field] = message
	}
}

func (v *Validator) Check(ok bool, field, message string) {
	if !ok {
		v.AddError(field, message)
	}
}

func (v *Validator) Required(value string, field string) {
	v.Check(strings.TrimSpace(value) != "", field, "This field is required")
}

func (v *Validator) MinLength(value string, minLength int, field string) {
	v.Check(len(value) >= minLength, field, fmt.Sprintf("Must be at least %d characters", minLength))
}

func (v *Validator) MaxLength(value string, maxLength int, field string) {
	v.Check(len(value) <= maxLength, field, fmt.Sprintf("Must be at most %d characters", maxLength))
}

func (v *Validator) Email(value string, field string) {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	v.Check(emailRegex.MatchString(value), field, "Must be a valid email address")
}

func (v *Validator) Username(value string, field string) {
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	v.Check(usernameRegex.MatchString(value), field, "Must be 3-20 characters, alphanumeric and underscore only")
}

func (v *Validator) Password(value string, field string) {
	v.MinLength(value, 8, field)

	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, char := range value {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsLower(char) {
			hasLower = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	v.Check(hasUpper, field, "Must contain at least one uppercase letter")
	v.Check(hasLower, field, "Must contain at least one lowercase letter")
	v.Check(hasDigit, field, "Must contain at least one digit")
}

func (v *Validator) In(value string, validValues []string, field string) {
	for _, valid := range validValues {
		if value == valid {
			return
		}
	}
	v.AddError(field, fmt.Sprintf("Must be one of: %s", strings.Join(validValues, ", ")))
}

func (v *Validator) MinValue(value, min int, field string) {
	v.Check(value >= min, field, fmt.Sprintf("Must be at least %d", min))
}

func (v *Validator) MaxValue(value, max int, field string) {
	v.Check(value <= max, field, fmt.Sprintf("Must be at most %d", max))
}

func (v *Validator) Range(value, min, max int, field string) {
	v.Check(value >= min && value <= max, field, fmt.Sprintf("Must be between %d and %d", min, max))
}

func (v *Validator) URL(value string, field string) {
	urlRegex := regexp.MustCompile(`^https?://[^\s]+$`)
	v.Check(urlRegex.MatchString(value), field, "Must be a valid URL")
}

func (v *Validator) Hashtag(value string, field string) {
	hashtagRegex := regexp.MustCompile(`^[a-z0-9_]{2,50}$`)
	normalized := strings.ToLower(strings.TrimPrefix(value, "#"))
	v.Check(hashtagRegex.MatchString(normalized), field, "Must be 2-50 characters, alphanumeric and underscore only")
}

func (v *Validator) NotEmpty(slice []string, field string) {
	v.Check(len(slice) > 0, field, "Must contain at least one item")
}

func (v *Validator) UniqueStrings(slice []string, field string) {
	seen := make(map[string]bool)
	for _, item := range slice {
		if seen[item] {
			v.AddError(field, "Must contain unique values")
			return
		}
		seen[item] = true
	}
}

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsValidUsername(username string) bool {
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	return usernameRegex.MatchString(username)
}

func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsLower(char) {
			hasLower = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}
