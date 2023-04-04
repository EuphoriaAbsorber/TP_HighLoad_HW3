package delivery

import (
	"dbproject/model"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// CreateUser godoc
// @Summary Creates User
// @Description Creates User
// @ID CreateUser
// @Accept  json
// @Produce  json
// @Tags User
// @Param nickname path string true "nickname of user"
// @Param user body model.User true "User params"
// @Success 201 {object} model.Response "OK"
// @Failure 400 {object} model.Error "Bad request - Problem with the request"
// @Failure 409 {object} model.Error "Conflict - User already exists"
// @Failure 500 {object} model.Error "Internal Server Error - Request is valid but operation failed at server side"
// @Router /user/{nickname}/create [post]
func (api *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	nickname := s[len(s)-2]
	decoder := json.NewDecoder(r.Body)
	var req model.User
	err := decoder.Decode(&req)
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, model.ErrBadRequest400, 400)
		return
	}
	req.Nickname = nickname

	users, err := api.usecase.GetUsersByUsermodel(&req)
	if err != nil {
		log.Println("get GetUserByUsermodel ", err)
		ReturnErrorJSON(w, model.ErrServerError500, 500)
		return
	}

	if len(users) > 0 {
		w.WriteHeader(409)
		json.NewEncoder(w).Encode(&users)
		return
	}

	err = api.usecase.CreateUser(&req)
	if err != nil {
		log.Println("db err: ", err)
		ReturnErrorJSON(w, model.ErrServerError500, 500)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(&req)
}

// GetProfile godoc
// @Summary Gets Users profile
// @Description Gets Users profile
// @ID GetProfile
// @Accept  json
// @Produce  json
// @Tags User
// @Param nickname path string true "nickname of user"
// @Success 200 {object} model.User
// @Failure 404 {object} model.Error "Not found - Requested entity is not found in database"
// @Failure 500 {object} model.Error "Internal Server Error - Request is valid but operation failed at server side"
// @Router /user/{nickname}/profile [get]
func (api *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	nickname := s[len(s)-2]

	user, err := api.usecase.GetProfile(nickname)
	if err == model.ErrNotFound404 {
		log.Println(err)
		ReturnErrorJSON(w, model.ErrNotFound404, 404)
		return
	}
	if err != nil {
		log.Println(err)
		ReturnErrorJSON(w, model.ErrServerError500, 500)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&user)
}

// PostProfile godoc
// @Summary Changes Users profile
// @Description Changes Users profile
// @ID PostProfile
// @Accept  json
// @Produce  json
// @Tags User
// @Param nickname path string true "nickname of user"
// @Param user body model.User true "User params"
// @Success 200 {object} model.Response "OK"
// @Failure 404 {object} model.Error "Not found - Requested entity is not found in database"
// @Failure 409 {object} model.Error "Conflict - User already exists"
// @Failure 500 {object} model.Error "Internal Server Error - Request is valid but operation failed at server side"
// @Router /user/{nickname}/profile [post]
func (api *Handler) PostProfile(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	nickname := s[len(s)-2]

	decoder := json.NewDecoder(r.Body)
	var req model.User
	err := decoder.Decode(&req)
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, model.ErrBadRequest400, 400)
		return
	}
	req.Nickname = nickname

	user, err := api.usecase.GetProfile(nickname)
	if err == model.ErrNotFound404 {
		log.Println(err)
		ReturnErrorJSON(w, model.ErrNotFound404, 404)
		return
	}
	if err != nil {
		log.Println(err)
		ReturnErrorJSON(w, model.ErrServerError500, 500)
		return
	}

	if req.Email != user.Email {
		users, err := api.usecase.GetUsersByUsermodel(&model.User{Email: req.Email, Nickname: ""})
		if err != nil {
			log.Println(err)
			ReturnErrorJSON(w, model.ErrServerError500, 500)
			return
		}

		if len(users) > 0 {
			ReturnErrorJSON(w, model.ErrConflict409, 409)
			return
		}
	}
	if req.Email == "" {
		req.Email = user.Email
	}
	if req.Fullname == "" {
		req.Fullname = user.Fullname
	}
	if req.About == "" {
		req.About = user.About
	}
	err = api.usecase.ChangeProfile(&req)
	if err != nil {
		log.Println(err)
		ReturnErrorJSON(w, model.ErrServerError500, 500)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&req)
}
