package errors

import "fmt"

func UserNotFoundById() ServiceError {
	return Error400(fmt.Errorf("user not found"), map[string]string{
		"ru": "Пользователь не найден, неверный ID",
	})

}

func WrongCredentials() ServiceError {
	return Error401(fmt.Errorf("wrong credentials"), map[string]string{
		"ru": "Неверный логин или пароль",
	})

}

func WrongRequest() ServiceError {
	return Error400(fmt.Errorf("wrong request"), map[string]string{
		"ru": "Не найдены данные за этот месяц / год",
	})

}

func AlreadyRegistered() ServiceError {
	return Error401(fmt.Errorf("already registered"), map[string]string{
		"ru": "Пользователь с таким логином уже существует",
		"en": "A user with this login already exists",
	})

}
func CreateUserErr(err error) ServiceError {
	return Error500(err, map[string]string{
		"ru": "Ошибка при регистрации пользователя",
		"en": "Error while registration",
	})

}
