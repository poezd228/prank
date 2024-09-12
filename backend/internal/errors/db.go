package errors

func DatabaseError(err error) ServiceError {
	return Error500(err, map[string]string{
		"en": "Database error. Please, try again later",
		"ru": "Ошибка обработки запроса. Пожалуйста, попробуйте позже",
	})
}
