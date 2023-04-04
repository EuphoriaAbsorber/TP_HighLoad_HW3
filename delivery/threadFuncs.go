package delivery

import (
	"dbproject/model"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// CreatePosts godoc
// @Summary Creates Posts
// @Description Creates Posts
// @ID CreatePosts
// @Accept  json
// @Produce  json
// @Tags Thread
// @Param slug_or_id path string true "slug or id"
// @Param posts body model.Posts true "Posts params"
// @Success 201 {object} model.Response "OK"
// @Failure 400 {object} model.Error "Bad request - Problem with the request"
// @Failure 404 {object} model.Error "Not found - Requested entity is not found in database"
// @Failure 409 {object} model.Error "Conflict - User already exists"
// @Failure 500 {object} model.Error "Internal Server Error - Request is valid but operation failed at server side"
// @Router /thread/{slug_or_id}/create [post]
func (api *Handler) CreatePosts(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	slug_or_id := s[len(s)-2]
	id := 0
	slug := slug_or_id
	id, err := strconv.Atoi(slug_or_id)
	if err != nil {
		log.Println("error: ", err)
	}

	decoder := json.NewDecoder(r.Body)
	req := new(model.Posts)
	err = decoder.Decode(&req)
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, model.ErrBadRequest400, 400)
		return
	}
	posts, err := api.usecase.CreatePosts(req, id, slug)
	if err == model.ErrNotFound404 {
		log.Println(err)
		ReturnErrorJSON(w, model.ErrNotFound404, 404)
		return
	}
	if err == model.ErrConflict409 {
		log.Println(err)
		ReturnErrorJSON(w, model.ErrConflict409, 409)
		return
	}
	if err != nil {
		log.Println(err)
		ReturnErrorJSON(w, model.ErrServerError500, 500)
		return
	}
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(&posts)
}

// GetThreadInfo godoc
// @Summary Gets thread info
// @Description Gets thread info
// @ID GetThreadInfo
// @Accept  json
// @Produce  json
// @Tags Thread
// @Param slug_or_id path string true "slug or id of thread"
// @Success 200 {object} model.Thread
// @Failure 404 {object} model.Error "Not found - Requested entity is not found in database"
// @Failure 500 {object} model.Error "Internal Server Error - Request is valid but operation failed at server side"
// @Router /thread/{slug_or_id}/details [get]
func (api *Handler) GetThreadInfo(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	slug_or_id := s[len(s)-2]
	id := 0
	slug := slug_or_id
	id, err := strconv.Atoi(slug_or_id)
	if err != nil {
		log.Println("error: ", err)
	}

	thread, err := api.usecase.GetThread(id, slug)
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
	json.NewEncoder(w).Encode(&thread)
}

// UpdateThreadInfo godoc
// @Summary Updates thread info
// @Description Updates thread info
// @ID UpdateThreadInfo
// @Accept  json
// @Produce  json
// @Tags Thread
// @Param slug_or_id path string true "slug or id of thread"
// @Param threadUpdate body model.ThreadUpdate true "ThreadUpdate params"
// @Success 200 {object} model.Thread
// @Failure 404 {object} model.Error "Not found - Requested entity is not found in database"
// @Failure 500 {object} model.Error "Internal Server Error - Request is valid but operation failed at server side"
// @Router /thread/{slug_or_id}/details [post]
func (api *Handler) UpdateThreadInfo(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	slug_or_id := s[len(s)-2]
	id := 0
	slug := slug_or_id
	id, err := strconv.Atoi(slug_or_id)
	if err != nil {
		log.Println("error: ", err)
	}
	decoder := json.NewDecoder(r.Body)
	var req model.ThreadUpdate
	err = decoder.Decode(&req)
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, model.ErrBadRequest400, 400)
		return
	}
	thread, err := api.usecase.UpdateThreadInfo(&req, id, slug)
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
	json.NewEncoder(w).Encode(&thread)
}

// VoteForThread godoc
// @Summary VoteForThread
// @Description VoteForThread
// @ID VoteForThread
// @Accept  json
// @Produce  json
// @Tags Thread
// @Param slug_or_id path string true "slug or id of thread"
// @Param vote body model.Vote true "vote params"
// @Success 200 {object} model.Thread
// @Failure 404 {object} model.Error "Not found - Requested entity is not found in database"
// @Failure 500 {object} model.Error "Internal Server Error - Request is valid but operation failed at server side"
// @Router /thread/{slug_or_id}/vote [post]
func (api *Handler) VoteForThread(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	slug_or_id := s[len(s)-2]
	id := 0
	slug := slug_or_id
	id, err := strconv.Atoi(slug_or_id)
	if err != nil {
		log.Println("error: ", err)
	}
	decoder := json.NewDecoder(r.Body)
	var req model.Vote
	err = decoder.Decode(&req)
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, model.ErrBadRequest400, 400)
		return
	}
	thread, err := api.usecase.VoteForThread(&req, id, slug)
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
	json.NewEncoder(w).Encode(&thread)
}

// GetThreadPosts godoc
// @Summary GetThreadPosts
// @Description GetThreadPosts
// @ID GetThreadPosts
// @Accept  json
// @Produce  json
// @Tags Thread
// @Param slug_or_id path string true "slug or id of thread"
// @Param   limit   query     string  false  "limit"
// @Param   since   query     string  false  "since"
// @Param   sort   query     string  false  "sort"
// @Param   desc    query     bool  false  "desc"
// @Success 200 {object} model.Threads
// @Failure 404 {object} model.Error "Not found - Requested entity is not found in database"
// @Failure 500 {object} model.Error "Internal Server Error - Request is valid but operation failed at server side"
// @Router /thread/{slug_or_id}/posts [get]
func (api *Handler) GetThreadPosts(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	slug_or_id := s[len(s)-2]
	id := 0
	slug := slug_or_id
	id, err := strconv.Atoi(slug_or_id)
	if err != nil {
		log.Println("error: ", err)
	}
	sinceS := r.URL.Query().Get("since")
	sort := r.URL.Query().Get("sort")
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
	since, err := strconv.Atoi(sinceS)
	if err != nil {
		log.Println("error: ", err)
		since = 0
	}
	if sort == "" {
		sort = "flat"
	}

	posts, err := api.usecase.GetThreadPosts(slug, id, limit, since, sort, desc)
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
	json.NewEncoder(w).Encode(&posts)
}
