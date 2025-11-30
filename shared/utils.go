package shared

import (
	"regexp"
	"strconv"
	"strings"
)

// SanitizeInput removes potentially dangerous characters
func SanitizeInput(input string) string {
	// Remove SQL injection attempts
	dangerous := []string{"--", ";", "/*", "*/", "xp_", "sp_", "DROP", "DELETE", "INSERT", "UPDATE", "EXEC"}
	result := input
	for _, d := range dangerous {
		result = strings.ReplaceAll(result, d, "")
	}
	return strings.TrimSpace(result)
}

// ValidatePrice validates price is a positive number
func ValidatePrice(price interface{}) (int, error) {
	switch v := price.(type) {
	case int:
		if v < 0 {
			return 0, NewInvalidInputError("Harga harus positif")
		}
		return v, nil
	case float64:
		if v < 0 {
			return 0, NewInvalidInputError("Harga harus positif")
		}
		return int(v), nil
	case string:
		p, err := strconv.Atoi(v)
		if err != nil {
			return 0, NewInvalidInputError("Harga harus berupa angka")
		}
		if p < 0 {
			return 0, NewInvalidInputError("Harga harus positif")
		}
		return p, nil
	default:
		return 0, NewInvalidInputError("Format harga tidak valid")
	}
}

// ValidateNotEmpty validates string is not empty
func ValidateNotEmpty(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return NewInvalidInputError(fieldName + " tidak boleh kosong")
	}
	return nil
}

// ValidatePhotoURL validates photo URL format
func ValidatePhotoURL(url string) error {
	if url == "" {
		return nil // optional field
	}

	pattern := `^https?://.*\.(jpg|jpeg|png|gif|webp)$`
	matched, err := regexp.MatchString(pattern, strings.ToLower(url))
	if err != nil || !matched {
		return NewInvalidInputError("URL foto tidak valid")
	}
	return nil
}

// FormatPrice formats price with Rp currency
func FormatPrice(price int) string {
	return "Rp " + strconv.Itoa(price)
}

// Contains checks if slice contains value
func Contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// RemoveFromSlice removes value from slice
func RemoveFromSlice(slice []string, value string) []string {
	result := []string{}
	for _, v := range slice {
		if v != value {
			result = append(result, v)
		}
	}
	return result
}
