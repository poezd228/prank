package validators

import (
	"regexp"

	"github.com/mvd-inc/anibliss/internal/errors"
)

func ValidateUserPasswordAndLogin(login string, password string) errors.ServiceError {
	if len(login) < 4 || len(login) > 15 {
		return errors.Error401(nil, map[string]string{
			"ru": "Логин должен содержать от 4 до 15 символов",
		})
	}

	loginRegex := regexp.MustCompile(`^[a-zA-Z0-9_ ]+$`)
	if !loginRegex.MatchString(login) {
		return errors.Error401(nil, map[string]string{
			"ru": "Логин может содержать только латиницу, цифры, пробелы и символы подчеркивания",
		})
	}

	// Валидация пароля
	if len(password) < 8 {
		return errors.Error401(nil, map[string]string{
			"ru": "Пароль должен содержать минимум 8 символов",
		})
	}

	passwordRegex := regexp.MustCompile(`^[a-zA-Z0-9!@#\$%\^&\*\(\)_\+\-=\[\]\{\};':"\\|,.<>\/?]+$`)
	if !passwordRegex.MatchString(password) {
		return errors.Error401(nil, map[string]string{
			"ru": "Пароль может содержать только латиницу, цифры и спец. символы",
		})
	}

	if !regexp.MustCompile(`\d`).MatchString(password) {
		return errors.Error401(nil, map[string]string{
			"ru": "Пароль должен содержать хотя бы одну цифру",
		})
	}

	return nil
}
