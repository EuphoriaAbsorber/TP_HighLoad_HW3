package delivery

import (
	"dbproject/model"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// GetPostById godoc
// @Summary Gets post by id
// @Description Gets post by id
// @ID GetPostById
// @Accept  json
// @Produce  json
// @Tags Post
// @Param id path string true "id of post"
// @Param   related   query     string  false  "related"
// @Success 200 {object} model.PostFull
// @Failure 404 {object} model.Error "Not found - Requested entity is not found in database"
// @Failure 500 {object} model.Error "Internal Server Error - Request is valid but operation failed at server side"
// @Router /post/{id}/details [get]
func (api *Handler) GetPostById(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	idS := s[len(s)-2]
	id, err := strconv.Atoi(idS)
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, model.ErrBadRequest400, 400)
		return
	}
	userFlag := false
	forumFlag := false
	threadFlag := false
	relatedS := r.URL.Query().Get("related")
	relatedArr := strings.Split(relatedS, ",")
	for _, rel := range relatedArr {
		if rel == "user" {
			userFlag = true
		}
		if rel == "forum" {
			forumFlag = true
		}
		if rel == "thread" {
			threadFlag = true
		}
	}
	answer := model.PostFull{}
	post, err := api.usecase.GetPostById(id)
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
	answer.Post = post
	if userFlag {
		user, err := api.usecase.GetProfile(post.Author)
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
		answer.Author = user
	}
	if forumFlag {
		forum, err := api.usecase.GetForumBySlug(post.Forum)
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
		answer.Forum = forum
	}
	if threadFlag {
		thread, err := api.usecase.GetThread(post.Thread, "")
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
		answer.Thread = thread
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&answer)
}

// UpdatePost godoc
// @Summary Editss post by id
// @Description Edits post by id
// @ID UpdatePost
// @Accept  json
// @Produce  json
// @Tags Post
// @Param id path string true "id of post"
// @Param message body model.PostUpdate true "PostUpdate params"
// @Success 200 {object} model.PostFull
// @Failure 404 {object} model.Error "Not found - Requested entity is not found in database"
// @Failure 500 {object} model.Error "Internal Server Error - Request is valid but operation failed at server side"
// @Router /post/{id}/details [post]
func (api *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	idS := s[len(s)-2]
	id, err := strconv.Atoi(idS)
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, model.ErrBadRequest400, 400)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var req model.PostUpdate
	err = decoder.Decode(&req)
	if err != nil {
		log.Println("error: ", err)
		ReturnErrorJSON(w, model.ErrBadRequest400, 400)
		return
	}

	post, err := api.usecase.UpdatePost(id, req.Message)
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
	json.NewEncoder(w).Encode(&post)
}
