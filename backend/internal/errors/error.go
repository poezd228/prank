package errors

import (
	"errors"
	"fmt"
)

type ServiceError interface {
	Error() string
	GetCode() int
	GetReason(lang string) string
}
type ServiceErr struct {
	Code    int
	Reasons map[string]string
	Err     error
}

func ErrorWithCode(code int, err error, reasons map[string]string) ServiceError {
	return ServiceErr{
		Code:    code,
		Err:     err,
		Reasons: reasons,
	}
}
func (err ServiceErr) Error() string {
	return err.Err.Error()
}

func (err ServiceErr) GetCode() int {
	return err.Code
}

func (err ServiceErr) GetReason(lang string) string {
	return err.Reasons[lang]
}
func Error500(err error, reasons map[string]string) ServiceError {
	return ErrorWithCode(500, err, reasons)
}

func Error503(err error, reasons map[string]string) ServiceError {
	return ErrorWithCode(503, err, reasons)
}

func Error404(err error, reasons map[string]string) ServiceError {
	return ErrorWithCode(404, err, reasons)
}

func Error403(err error, reasons map[string]string) ServiceError {
	return ErrorWithCode(403, err, reasons)
}

func Error400(err error, reasons map[string]string) ServiceError {
	return ErrorWithCode(400, err, reasons)
}

func Error401(err error, reasons map[string]string) ServiceError {
	return ErrorWithCode(401, err, reasons)
}

func InternalError(err error) ServiceError {
	return Error500(err, map[string]string{
		"en": "Internal error, please, try again later",
		"ru": "Внутренняя ошибка. Пожалуйста, попробуйте снова позже",
	})
}

func ParseFailed(err error) ServiceError {
	return Error400(err, map[string]string{
		"en": "Parse failed. Please, contact support",
		"ru": "Ошибка обработки запроса. Пожалуйста, обратитесь в поддержку",
	})
}

func ValidationFailed(err error) ServiceError {
	return Error400(err, map[string]string{
		"en": "Validation failed. Please, contact support",
		"ru": "Ошибка валидации запроса. Пожалуйста, обратитесь в поддержку",
	})
}

func DisabledDebugMode() ServiceError {
	return Error400(fmt.Errorf("debug mode disabled"), map[string]string{
		"en": "Debug mode disabled",
		"ru": "Режим отладки заблокирован",
	})
}

var (
	AuthAuthFailedRaw = errors.New("auth failed")
	AuthAuthFailed    = &ServiceErr{
		Code: 401,
		Reasons: map[string]string{
			"en": "auth failed",
		},
		Err: AuthAuthFailedRaw,
	}

	AuthParseTokenRaw = errors.New("parse token failed")
	AuthParseToken    = &ServiceErr{
		Code: 400,
		Reasons: map[string]string{
			"en": AuthParseTokenRaw.Error(),
		},
		Err: AuthParseTokenRaw,
	}
	InvalidTokenClaims = &ServiceErr{
		Code: 400,
		Reasons: map[string]string{
			"en": "invalid token claims",
		},
		Err: errors.New("invalid token claims"),
	}

	AuthHashPassword = &ServiceErr{
		Code: 400,
		Reasons: map[string]string{
			"en": "hashing password error",
		},
		Err: errors.New("hashing password error"),
	}

	AuthExpiredToken = &ServiceErr{
		Code: 400,
		Reasons: map[string]string{
			"en": "expired token",
		},
		Err: errors.New("expired token"),
	}

	AuthInvalidTokenPurpose = &ServiceErr{
		Code: 400,
		Reasons: map[string]string{
			"en": "invalid token purpose",
		},
		Err: errors.New("invalid token purpose"),
	}

	AuthInvalidToken = &ServiceErr{
		Code: 400,
		Reasons: map[string]string{
			"en": "invalid token",
		},
		Err: errors.New("invalid token"),
	}

	AuthCreateTokens = &ServiceErr{
		Code: 400,
		Reasons: map[string]string{
			"en": "create tokens error",
		},
		Err: errors.New("create tokens error"),
	}

	AuthUserNotFoundByIdRaw = errors.New("there is no user with such id")
	AuthUserNotFoundById    = &ServiceErr{
		Code: 400,
		Reasons: map[string]string{
			"en": AuthUserNotFoundByIdRaw.Error(),
		},
		Err: AuthUserNotFoundByIdRaw,
	}

	AuthNumberAssignmentFailedRaw = errors.New("number assignment failed")
	AuthNumberAssignmentFailed    = &ServiceErr{
		Code: 400,
		Reasons: map[string]string{
			"en": AuthNumberAssignmentFailedRaw.Error(),
		},
		Err: AuthNumberAssignmentFailedRaw,
	}

	AuthGetUserDataFailed = &ServiceErr{
		Code: 400,
		Reasons: map[string]string{
			"en": "get user data failed",
		},
		Err: errors.New("get user data failed"),
	}

	CreateToken = &ServiceErr{
		Code: 400,
		Reasons: map[string]string{
			"en": "token not created",
		},
		Err: errors.New("token not created"),
	}

	TokenDoesNotExist = &ServiceErr{
		Code: 400,
		Reasons: map[string]string{
			"en": "token does not exist",
		},
		Err: errors.New("token does not exist"),
	}
)
