package dto

import (
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

type SignUpRequest struct {
	Name            string `json:"name,omitempty"`
	Email           string `json:"email,omitempty"`
	Password        string `json:"password,omitempty"`
	ConfirmPassword string `json:"confirm_password,omitempty"`
}

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^[6-9]\d{9}$`)
	zipRegex   = regexp.MustCompile(`^[1-9]\d{5}$`)
)

func (r *SignUpRequest) Validate() map[string][]string {
	errs := make(map[string][]string)

	if strings.TrimSpace(r.Name) == "" {
		errs["name"] = append(errs["name"], "name is required")
	}

	if strings.TrimSpace(r.Email) == "" {
		errs["email"] = append(errs["email"], "email is required")
	} else if !emailRegex.MatchString(r.Email) {
		errs["email"] = append(errs["email"], "email is invalid")
	}

	validatePassword(r.Password, errs)

	if strings.TrimSpace(r.ConfirmPassword) == "" {
		errs["confirm_password"] = append(errs["confirm_password"], "confirm password is required")
	} else if r.Password != r.ConfirmPassword {
		errs["confirm_password"] = append(errs["confirm_password"], "passwords do not match")
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

func validatePassword(password string, errs map[string][]string) {
	if strings.TrimSpace(password) == "" {
		errs["password"] = append(errs["password"], "password is required")
	} else {
		if len(password) < 8 {
			errs["password"] = append(errs["password"], "password must be at least 8 characters")
		}

		var hasUpper, hasLower, hasDigit, hasSpecial bool

		for _, c := range password {
			switch {
			case unicode.IsUpper(c):
				hasUpper = true
			case unicode.IsLower(c):
				hasLower = true
			case unicode.IsDigit(c):
				hasDigit = true
			case unicode.IsPunct(c) || unicode.IsSymbol(c):
				hasSpecial = true
			}
		}

		if !hasUpper {
			errs["password"] = append(errs["password"], "password must contain an uppercase letter")
		}

		if !hasLower {
			errs["password"] = append(errs["password"], "password must contain a lowercase letter")
		}

		if !hasDigit {
			errs["password"] = append(errs["password"], "password must contain a number")
		}

		if !hasSpecial {
			errs["password"] = append(errs["password"], "password must contain a special character")
		}
	}
}

type SignInRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (r *SignInRequest) Validate() map[string][]string {
	errs := make(map[string][]string)

	if strings.TrimSpace(r.Email) == "" {
		errs["email"] = append(errs["email"], "email is required")
	} else if !emailRegex.MatchString(r.Email) {
		errs["email"] = append(errs["email"], "email is invalid")
	}

	validatePassword(r.Password, errs)

	if len(errs) == 0 {
		return nil
	}

	return errs
}

type SignInResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type CreateAddressRequest struct {
	UserID    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  *string   `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address1  string    `json:"address1"`
	Address2  *string   `json:"address2"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	Zip       string    `json:"zip"`
}

func (r *CreateAddressRequest) Validate() map[string][]string {
	errs := make(map[string][]string)

	if strings.TrimSpace(r.FirstName) == "" {
		errs["first_name"] = append(errs["first_name"], "first name is required")
	}

	if r.LastName != nil && strings.TrimSpace(*r.LastName) == "" {
		errs["last_name"] = append(errs["last_name"], "last name cannot be empty")
	}

	if strings.TrimSpace(r.Email) == "" {
		errs["email"] = append(errs["email"], "email is required")
	} else if !emailRegex.MatchString(r.Email) {
		errs["email"] = append(errs["email"], "email is invalid")
	}

	if strings.TrimSpace(r.Phone) == "" {
		errs["phone"] = append(errs["phone"], "phone number is required")
	} else if !phoneRegex.MatchString(r.Phone) {
		errs["phone"] = append(errs["phone"], "phone number is invalid")
	}

	if strings.TrimSpace(r.Address1) == "" {
		errs["address1"] = append(errs["address1"], "address line 1 is required")
	}

	if r.Address2 != nil && strings.TrimSpace(*r.Address2) == "" {
		errs["address2"] = append(errs["address2"], "address line 2 cannot be empty")
	}

	if strings.TrimSpace(r.City) == "" {
		errs["city"] = append(errs["city"], "city is required")
	}

	if strings.TrimSpace(r.State) == "" {
		errs["state"] = append(errs["state"], "state is required")
	}

	if strings.TrimSpace(r.Zip) == "" {
		errs["zip"] = append(errs["zip"], "zip code is required")
	} else if !zipRegex.MatchString(r.Zip) {
		errs["zip"] = append(errs["zip"], "zip code is invalid")
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

type AddressResponse struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	FirstName string
	LastName  *string
	Email     string
	Phone     string
	Address1  string
	Address2  *string
	City      string
	State     string
	Zip       string

	CreatedAt time.Time
	UpdatedAt time.Time
}
