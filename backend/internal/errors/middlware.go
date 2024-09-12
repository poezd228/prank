package errors

func MissingCookies(err error) ServiceError {
	return Error400(err, map[string]string{
		"ru": "Отсутствуют куки",
		"en": "Cookies are missing",
	})

}
