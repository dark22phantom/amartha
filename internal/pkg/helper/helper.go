package helper

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateAccessToken(secretKey, adminID, adminName string) (string, error) {
	claims := jwt.MapClaims{
		"admin_id":   adminID,
		"admin_name": adminName,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

var AllowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".pdf":  true,
}

func ValidateFileExtension(filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	if _, ok := AllowedExtensions[ext]; !ok {
		return fmt.Errorf("invalid file extension: %s", ext)
	}
	return nil
}
