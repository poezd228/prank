package writers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mvd-inc/anibliss/internal/domain"
	"github.com/mvd-inc/anibliss/internal/errors"
	"github.com/mvd-inc/anibliss/models"
	"github.com/mvd-inc/anibliss/pkg/logger"
)

func WriteResponse(w http.ResponseWriter, code int64, resp any) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(code))

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(resp); err != nil {
		return fmt.Errorf("encode write resp: %w", err)
	}

	return nil
}

func WriteErrorResponseWithErrorLog(w http.ResponseWriter, l logger.Logger, e errors.ServiceError, lang string) {
	err := WriteErrorResponse(w, e, lang)
	if err != nil {
		l.Error("write error response failed", err)
	}
}

func WriteErrorResponse(w http.ResponseWriter, err errors.ServiceError, lang string) error {
	code := int64(err.GetCode())
	message := err.GetReason(lang)
	enMessage := err.GetReason("en")
	detail := err.Error()
	return WriteResponse(w, code, models.ErrorResponse{
		Code:      code,
		Detail:    detail,
		Message:   message,
		EnMessage: enMessage,
	})
}

func WriteResponseWithErrorLog(w http.ResponseWriter, l logger.Logger, code int64, resp any) {
	err := WriteResponse(w, code, resp)
	if err != nil {
		l.Error("write response failed", err)
	}
}

func WriteTokenResponseWithErrorLog(w http.ResponseWriter, l logger.Logger, code int64, tokens domain.TokenResponse) {

	err := WriteTokenResponse(w, code, tokens)
	if err != nil {
		l.Error("write tokens failed", err)
	}

}
func WriteChangePasswordResponseWithErrorLog(w http.ResponseWriter, l logger.Logger, code int64) {
	encoder := json.NewEncoder(w)
	message := "Пароль успешно изменен"
	enMassage := "Password successfully changed"

	encoder.SetEscapeHTML(false)
	resp := &models.ChangePasswordResponse{
		Code:      code,
		Message:   message,
		EnMessage: enMassage,
		Detail:    "",
	}
	err := WriteResponse(w, code, resp)
	if err != nil {
		l.Error("write change pass response failed")
	}

}

func WriteTokenResponse(w http.ResponseWriter, code int64, tokens domain.TokenResponse) error {

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	accessTokenCookie := &http.Cookie{
		Name:     "access-token",
		Value:    tokens.AccessToken.Token,
		Expires:  time.Unix(0, tokens.AccessToken.ExpiresAt*1e6),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	http.SetCookie(w, accessTokenCookie)

	refreshTokenCookie := &http.Cookie{
		Name:     "refresh-token",
		Value:    tokens.RefreshToken.Token,
		Expires:  time.Unix(0, tokens.RefreshToken.ExpiresAt*1e6),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	http.SetCookie(w, refreshTokenCookie)
	message := "Успешная авторизация"
	enMessage := "Successfully auth"
	resp := models.AuthResponse{
		Code:      code,
		Message:   message,
		EnMessage: enMessage,
		Detail:    "",
	}
	err := WriteResponse(w, code, resp)
	if err != nil {
		return err
	}
	return nil
}
