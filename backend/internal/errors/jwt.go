package errors

func ParseTokenFailed(err error) ServiceError {
	return Error401(err, map[string]string{
		"en": "Parse token failed. Please, contact support",
		"ru": "Ошибка обработки запроса. Пожалуйста, обратитесь в поддержку",
	})
}

func CheckTokenFailed(err error) ServiceError {
	return Error401(err, map[string]string{
		"en": "Invalid token",
		"ru": "Недействительный токен",
	})

}

func TokenAddFailed(err error) ServiceError {
	return Error400(err, map[string]string{
		"ru": "Ошибка добавления токена",
		"en": "Token add failed",
	})

}
func SignStringErr(err error) ServiceError {
	return Error400(err, map[string]string{
		"ru": "Ошибка подписи токена",
		"en": "Sign string error",
	})

}
func FindNumberError(err error) ServiceError {
	return Error400(err, map[string]string{
		"en": "finding number error",
	})

}
func DropTokensError(err error) ServiceError {
	return Error400(err, map[string]string{
		"en": "drop tokens error",
	})
}

func DropOldTokensError(err error) ServiceError {
	return Error400(err, map[string]string{
		"en": "drop old tokens error",
	})

}
