package v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/wanomir/e"
	"github.com/wanomir/rr"
	"net/http"
	"net/url"
	"strconv"
	"userservice/internal/entity"
)

type Servicer interface {
	GetUsers(ctx context.Context, offset, limit int) ([]entity.User, error)
	GetUser(ctx context.Context, email string) (entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (int, error)
}

type Controller struct {
	service Servicer
	rr      *rr.ReadResponder
}

func NewController(service Servicer) *Controller {
	return &Controller{
		service: service,
		rr:      rr.NewReadResponder(),
	}
}

// GetUsers godoc
// @Summary Returns a list of users provided offset and limit
// @Tags users
// @Produce json
// @Param offset query int true "offset"
// @Param limit query int true "limit"
// @Success 200 {object} rr.JSONResponse
// @Failure 400,500 {object} rr.JSONResponse
// @Router /users [get]
func (c *Controller) GetUsers(w http.ResponseWriter, r *http.Request) {
	requestUrl, _ := url.Parse(r.URL.String())

	limit, err := strconv.Atoi(requestUrl.Query().Get("limit"))
	if err != nil {
		_ = c.rr.WriteJSONError(w, e.WrapIfErr("error parsing limit", err))
		return
	}

	offset, err := strconv.Atoi(requestUrl.Query().Get("offset"))
	if err != nil {
		_ = c.rr.WriteJSONError(w, e.WrapIfErr("error parsing offset", err))
		return
	}

	users, err := c.service.GetUsers(r.Context(), offset, limit)
	if err != nil {
		_ = c.rr.WriteJSONError(w, err, 500)
		return
	}

	resp := rr.JSONResponse{
		Message: fmt.Sprintf("%d users found", len(users)),
		Data:    users,
	}

	_ = c.rr.WriteJSON(w, 200, resp)
}

// GetUser godoc
// @Summary Returns user object provided user email
// @Tags users
// @Produce json
// @Param email query string true "email"
// @Success 200 {object} rr.JSONResponse
// @Failure 400,404 {object} rr.JSONResponse
// @Router /users/{email} [get]
func (c *Controller) GetUser(w http.ResponseWriter, r *http.Request) {
	requestUrl, _ := url.Parse(r.URL.String())

	email := requestUrl.Query().Get("email")
	if email == "" {
		_ = c.rr.WriteJSONError(w, errors.New("invalid email"))
		return
	}

	user, err := c.service.GetUser(r.Context(), email)
	if err != nil {
		_ = c.rr.WriteJSONError(w, err, 404)
		return
	}

	resp := rr.JSONResponse{
		Message: fmt.Sprintf("user with email %s found", email),
		Data:    user,
	}

	_ = c.rr.WriteJSON(w, 200, resp)
}

// CreateUser godoc
// @Summary Creates a new user provided user object
// @Tags users
// @Accept json
// @Produce json
// @Param user body entity.User true "user object"
// @Success 201 {object} rr.JSONResponse
// @Failure 400 {object} rr.JSONResponse
// @Router /users [post]
func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	if err := c.rr.ReadJSON(w, r, &user); err != nil {
		_ = c.rr.WriteJSONError(w, e.Wrap("error reading user", err))
		return
	}

	userId, err := c.service.CreateUser(r.Context(), user)
	if err != nil {
		_ = c.rr.WriteJSONError(w, e.Wrap("error creating user", err))
		return
	}

	resp := rr.JSONResponse{
		Message: fmt.Sprintf("created user with id %d", userId),
	}

	_ = c.rr.WriteJSON(w, 201, resp)
}
