package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"user-service/internal/store"

	"github.com/lib/pq"
)

var successMsg = map[string]bool{
	"success": true,
}

var (
	errorUnableToLogin  = errors.New("unable to login")
	errorInternalServer = errors.New("server error")
	errorInvalidId      = errors.New("invalid id")
	errorEmailTaken     = errors.New("this email has already been taken")
)

type apiResponse struct {
	Err    string `json:"error"`
	Data   any    `json:"data"`
	status int
}

type userCreateRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type apiHandler func(http.ResponseWriter, *http.Request) *apiResponse

// godoc
// @Summary		 wraps the apiHandler function
func makeHTTPHandler(f apiHandler, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		apiResp := f(w, r)
		if apiResp.Err != "" {
			writeJSON(w, apiResp, apiResp.status)
			return
		}
		writeJSON(w, apiResp, apiResp.status)
	}
}

func (a *API) registerUser(w http.ResponseWriter, r *http.Request) *apiResponse {
	resp := &apiResponse{}
	var userCrt userCreateRequest
	if err := readJSON(r, &userCrt); err != nil {
		resp.Err = err.Error()
		resp.status = http.StatusBadRequest
		return resp
	}

	hashedPass, err := hashPassword(userCrt.Password)
	if err != nil {
		log.Println("error hasing the password: ", err)
		resp.status = http.StatusInternalServerError
		resp.Err = err.Error()
		return resp
	}

	u := &store.User{
		Name:     userCrt.Name,
		LastName: userCrt.LastName,
		Password: hashedPass,
		Email:    userCrt.Email,
		Role:     store.CUSTOMER,
	}

	if err := a.store.Create(u); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == store.UniqueError {
			resp.status = http.StatusBadRequest
			resp.Err = errorEmailTaken.Error()
			return resp
		}
		log.Println("error quering the database: ", err)
		resp.Err = err.Error()
		resp.status = http.StatusInternalServerError
		return resp
	}

	resp.Data = successMsg
	resp.status = http.StatusOK
	return resp
}

func (a *API) login(w http.ResponseWriter, r *http.Request) *apiResponse {
	var login loginRequest
	resp := &apiResponse{}
	if err := readJSON(r, &login); err != nil {
		log.Println("error reading from request body", err)
		resp.status = http.StatusBadRequest
		resp.Err = err.Error()
		return resp
	}

	user, err := a.store.GetByEmail(login.Email)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			resp.status = http.StatusNotFound
			resp.Err = err.Error()
		default:
			log.Println("error getting user by email: ", err)
			resp.status = http.StatusInternalServerError
			resp.Err = errorInternalServer.Error()
		}
		return resp
	}

	if isValid := compareHash(login.Password, user.Password); !isValid {
		log.Println("error comparing password: ", err)
		resp.status = http.StatusBadRequest
		resp.Err = errorUnableToLogin.Error()
		return resp
	}

	resp.status = http.StatusOK
	resp.Data = successMsg
	return resp
}

func (a *API) deleteAccount(w http.ResponseWriter, r *http.Request) *apiResponse {
	id := r.PathValue("id")
	validId, err := strconv.Atoi(id)
	var resp = &apiResponse{}
	if err != nil {
		resp.status = http.StatusBadRequest
		resp.Err = errorInvalidId.Error()
		return resp
	}
	if err := a.store.Delete(validId); err != nil {
		switch err {
		case sql.ErrNoRows:
			resp.status = http.StatusNotFound
			resp.Err = errorInvalidId.Error()
		default:
			log.Println("error deleting a user: ", err)
			resp.status = http.StatusInternalServerError
			resp.Err = errorInternalServer.Error()
		}
		return resp

	}
	resp.status = http.StatusOK
	resp.Data = successMsg
	return resp
}

func (a *API) updateAccount(w http.ResponseWriter, r *http.Request) *apiResponse {
	resp := &apiResponse{}
	id := r.PathValue("id")
	validId, err := strconv.Atoi(id)
	if err != nil {
		resp.status = http.StatusBadRequest
		resp.Err = errorInvalidId.Error()
		return resp
	}
	var u userCreateRequest
	if err := readJSON(r, &u); err != nil {
		resp.status = http.StatusBadRequest
		resp.Err = err.Error()
		return resp
	}

	newUsr := &store.User{
		Name:     u.Name,
		LastName: u.LastName,
	}

	if err := a.store.Update(newUsr, validId); err != nil {
		resp.status = http.StatusInternalServerError
		resp.Err = errorInternalServer.Error()
		return resp
	}

	resp.status = http.StatusOK
	resp.Data = successMsg
	return resp
}
