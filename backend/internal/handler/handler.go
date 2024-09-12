package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mvd-inc/anibliss/internal/config"
	"github.com/mvd-inc/anibliss/internal/domain"
	"github.com/mvd-inc/anibliss/internal/errors"
	"github.com/mvd-inc/anibliss/internal/handler/middleware"
	"github.com/mvd-inc/anibliss/internal/handler/writers"
	"github.com/mvd-inc/anibliss/internal/service/auth"
	"github.com/mvd-inc/anibliss/internal/service/users"
	"github.com/mvd-inc/anibliss/models"
	"github.com/mvd-inc/anibliss/pkg/logger"
	"github.com/mvd-inc/anibliss/pkg/utils"
)

type Handler interface {
	FillHandlers(mux *http.ServeMux)
}

type handler struct {
	l            logger.Logger
	cfg          config.HandlerConfig
	middleware   middleware.Middleware
	usersService users.Service
	authService  auth.Service
}

func NewHandler(
	cfg config.HandlerConfig,
	middleware middleware.Middleware,
	l logger.Logger,
	usersService users.Service,
	authService auth.Service,
) Handler {
	return &handler{
		cfg:          cfg,
		middleware:   middleware,
		l:            l,
		usersService: usersService,
		authService:  authService,
	}
}

func addHandlerToMux(mux *http.ServeMux, method string, base string, path string, handler http.Handler) {
	mux.Handle(fmt.Sprintf("%s %s", method, base+path), handler)
	mux.Handle(fmt.Sprintf("%s %s", http.MethodOptions, base+path), handler)
	mux.Handle(fmt.Sprintf("%s %s/", method, base+path), handler)
	mux.Handle(fmt.Sprintf("%s %s/", http.MethodOptions, base+path), handler)

}

func (h *handler) FillHandlers(mux *http.ServeMux) {
	// сдеся хендлеры в таком виде
	var base string
	base = "/auth"
	addHandlerToMux(mux, http.MethodPost, base, "/refresh", h.middleware.RefreshMiddleware(h.refresh))
	addHandlerToMux(mux, http.MethodPost, base, "/login", http.HandlerFunc(h.signIn))
	addHandlerToMux(mux, http.MethodPost, base, "/register", http.HandlerFunc(h.register))
	addHandlerToMux(mux, http.MethodPost, base, "/signin", h.middleware.AuthMiddleware(h.login))

}

func (h *handler) login(acc domain.Account, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	writers.WriteResponseWithErrorLog(w, h.l, http.StatusOK, acc)

}

func (h *handler) register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req models.SignInRequest
	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
	defer cancel()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writers.WriteErrorResponseWithErrorLog(w, h.l, errors.ParseFailed(err), "ru")
		return
	}
	e := h.usersService.CreateUser(ctx, req.Login, req.Password)
	if e != nil {
		writers.WriteErrorResponseWithErrorLog(w, h.l, e, "ru")
		return

	}
	writers.WriteResponseWithErrorLog(w, h.l, http.StatusOK, models.AuthResponse{
		Code:      200,
		Detail:    "Success",
		EnMessage: "Success",
		Message:   "Успешно",
	})

}
func (h *handler) changePassword(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var changePasswordReq models.ChangePassRequest
	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
	defer cancel()
	err := json.NewDecoder(r.Body).Decode(&changePasswordReq)
	if err != nil {
		writers.WriteErrorResponseWithErrorLog(w, h.l, errors.ValidationFailed(err), "ru")
		return
	}
	changePasswordReq.OldPass = utils.HashSha256(changePasswordReq.OldPass)
	changePasswordReq.NewPass = utils.HashSha256(changePasswordReq.NewPass)
	e := h.usersService.ChangeUserPass(ctx, changePasswordReq.OldPass, changePasswordReq.NewPass, changePasswordReq.Login)
	if e != nil {
		writers.WriteErrorResponseWithErrorLog(w, h.l, e, "ru")
		return
	}
	writers.WriteChangePasswordResponseWithErrorLog(w, h.l, http.StatusOK)

}

func (h *handler) signIn(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req models.SignInRequest
	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
	defer cancel()
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writers.WriteErrorResponseWithErrorLog(w, h.l, errors.ValidationFailed(err), "ru")
		return

	}
	acc, e := h.usersService.GetUserLP(ctx, req.Login, req.Password)
	if e != nil {
		writers.WriteErrorResponseWithErrorLog(w, h.l, e, "ru")
		return

	}
	resp, e := h.authService.AccAuth(ctx, acc)
	if e != nil {
		writers.WriteErrorResponseWithErrorLog(w, h.l, e, "ru")
		return
	}

	writers.WriteTokenResponseWithErrorLog(w, h.l, http.StatusOK, resp)

}
func (h *handler) refresh(acc domain.Account, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
	defer cancel()
	acc, err := h.usersService.GetUser(ctx, acc)
	if err != nil {
		writers.WriteErrorResponseWithErrorLog(w, h.l, err, "ru")
		return

	}
	resp := h.authService.JwtRefresh(ctx, acc)

	writers.WriteTokenResponseWithErrorLog(w, h.l, http.StatusOK, resp)

}
