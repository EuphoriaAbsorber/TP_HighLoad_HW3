package delivery

import (
	"dbproject/model"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// CreateForum godoc
// @Summary Creates Forum
// @Description Creates Forum
// @ID CreateForum
// @Accept  json
// @Produce  json
// @Tags Forum
// @Param forum body model.ForumCreateModel true "Forum params"
// @Success 201 {object} model.Response "OK"
// @Failure 400 {object} model.Error "Bad request - Problem with the request"
// @Failure 404 {object} model.Error "Not found - Requested entity is not found in database"
// @Failure 409 {object} model.Error "Conflict - User already exists"
// @Failure 500 {object} model.Error "Internal Server Error - Request is valid but operation failed at server side"
// @Router /forum/create [post]
func (api *Handler) CreateForum(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req model.Forum
	err := decoder.Decode(&req)
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, model.ErrBadRequest400, 400)
		return
	}
	user, err := api.usecase.GetProfile(req.User)
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
	req.User = user.Nickname

	forum, err := api.usecase.GetForumByUsername(req.User)
	if err != nil && err != model.ErrNotFound404 {
		log.Println(err)
		ReturnErrorJSON(w, model.ErrServerError500, 500)
		return
	}
	if forum != nil {
		w.WriteHeader(409)
		json.NewEncoder(w).Encode(&forum)
		return
	}
	err = api.usecase.CreateForum(&req)
	if err != nil {
		log.Println("db err: ", err)
		ReturnErrorJSON(w, model.ErrServerError500, 500)
		return
	}
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(&req)
}

// GetForumInfo godoc
// @Summary Gets forum info
// @Description Gets forum info
// @ID GetForumInfo
// @Accept  json
// @Produce  json
// @Tags Forum
// @Param slug path string true "slug of forum"
// @Success 200 {object} model.Forum
// @Failure 404 {object} model.Error "Not found - Requested entity is not found in database"
// @Failure 500 {object} model.Error "Internal Server Error - Request is valid but operation failed at server side"
// @Router /forum/{slug}/details [get]
func (api *Handler) GetForumInfo(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	slug := s[len(s)-2]

	forum, err := api.usecase.GetForumBySlug(slug)
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
	json.NewEncoder(w).Encode(&forum)
}

// CreateThread godoc
// @Summary creates thread
// @Description creates thread
// @ID CreateThread
// @Accept  json
// @Produce  json
// @Tags Forum
// @Param slug path string true "slug of forum"
// @Param thread body model.ThreadCreateModel true "Thread params"
// @Success 201 {object} model.Thread
// @Failure 404 {object} model.Error "Not found - Requested entity is not found in database"
// @Failure 409 {object} model.Thread
// @Failure 500 {object} model.Error "Internal Server Error - Request is valid but operation failed at server side"
// @Router /forum/{slug}/create [post]
func (api *Handler) CreateThread(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	slug := s[len(s)-2]
	var req model.Thread
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, model.ErrBadRequest400, 400)
		return
	}

	req.Votes = 0

	forum, err := api.usecase.GetForumBySlug(slug)
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
	req.Forum = forum.Slug
	user, err := api.usecase.GetProfile(req.Author)
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
	req.Author = user.Nickname
	thread, err := api.usecase.GetThread(0, req.Slug)
	if err != nil && err != model.ErrNotFound404 {
		log.Println(err)
		ReturnErrorJSON(w, model.ErrServerError500, 500)
		return
	}
	if thread != nil {
		w.WriteHeader(409)
		json.NewEncoder(w).Encode(&thread)
		return
	}
	thread, err = api.usecase.CreateThreadByModel(&req)
	if err != nil {
		log.Println(err)
		ReturnErrorJSON(w, model.ErrServerError500, 500)
		return
	}
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(&thread)
}

// GetForumUsers godoc
// @Summary Gets forum users
// @Description Gets forum users
// @ID GetForumUsers
// @Accept  json
// @Produce  json
// @Tags Forum
// @Param slug path string true "slug of forum"
// @Param   limit   query     string  false  "limit"
// @Param   since   query     string  false  "since"
// @Param   desc    query     bool  false  "desc"
// @Success 200 {object} model.Forum
// @Failure 404 {object} model.Error "Not found - Requested entity is not found in database"
// @Failure 500 {object} model.Error "Internal Server Error - Request is valid but operation failed at server side"
// @Router /forum/{slug}/users [get]
func (api *Handler) GetForumUsers(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	slug := s[len(s)-2]

	since := r.URL.Query().Get("since")
	limitS := r.URL.Query().Get("limit")
	descS := r.URL.Query().Get("desc")
	desc := false
	if descS == "true" {
		desc = true
	}
	limit, err := strconv.Atoi(limitS)
	if err != nil {
		log.Println("error: ", err)
		limit = 1e9
	}
	forum, err := api.usecase.GetForumBySlug(slug)
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

	users, err := api.usecase.GetForumUsers(forum.Slug, limit, since, desc)
	if err != nil {
		log.Println(err)
		ReturnErrorJSON(w, model.ErrServerError500, 500)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&users)
}

// GetForumThreads godoc
// @Summary Gets forum threads
// @Description Gets forum threads
// @ID GetForumThreads
// @Accept  json
// @Produce  json
// @Tags Forum
// @Param slug path string true "slug of forum"
// @Param   limit   query     string  false  "limit"
// @Param   since   query     string  false  "since"
// @Param   desc    query     bool  false  "desc"
// @Success 200 {object} model.Threads
// @Failure 404 {object} model.Error "Not found - Requested entity is not found in database"
// @Failure 500 {object} model.Error "Internal Server Error - Request is valid but operation failed at server side"
// @Router /forum/{slug}/threads [get]
func (api *Handler) GetForumThreads(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	slug := s[len(s)-2]

	sinceS := r.URL.Query().Get("since")
	limitS := r.URL.Query().Get("limit")
	descS := r.URL.Query().Get("desc")
	desc := false
	if descS == "true" {
		desc = true
	}
	limit, err := strconv.Atoi(limitS)
	if err != nil {
		log.Println("error: ", err)
		limit = 1e9
	}
	since, err := time.Parse(time.RFC3339, "1971-01-01T00:00:00.000Z")
	if err != nil {
		log.Println(err)
		ReturnErrorJSON(w, model.ErrBadRequest400, 400)
		return
	}
	if sinceS != "" {
		since, err = time.Parse(time.RFC3339, sinceS)
		if err != nil {
			log.Println(err)
			ReturnErrorJSON(w, model.ErrBadRequest400, 400)
			return
		}
	}

	forum, err := api.usecase.GetForumBySlug(slug)
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
	threads, err := api.usecase.GetForumThreads(forum.Slug, limit, since, desc)
	if err != nil {
		log.Println(err)
		ReturnErrorJSON(w, model.ErrServerError500, 500)
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&threads)
}
