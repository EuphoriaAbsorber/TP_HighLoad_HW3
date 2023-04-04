package delivery

import (
	"dbproject/model"
	usecase "dbproject/usecase"
	"encoding/json"
	"log"
	"net/http"
)

// @title DB project API
// @version 1.0
// @description DB project server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:5000
// @BasePath  /api

type Handler struct {
	usecase usecase.UsecaseInterface
}

func NewHandler(uc usecase.UsecaseInterface) *Handler {
	return &Handler{
		usecase: uc,
	}
}

func ReturnErrorJSON(w http.ResponseWriter, err error, errCode int) {
	w.WriteHeader(errCode)
	json.NewEncoder(w).Encode(&model.Error{Error: err.Error()})
}

// GetServiceStatus godoc
// @Summary Gets Service info
// @Description Gets Service info
// @ID GetServiceStatus
// @Accept  json
// @Produce  json
// @Tags Service
// @Success 200 {object} model.Status
// @Failure 500 {object} model.Error "Internal Server Error - Request is valid but operation failed at server side"
// @Router /service/status [get]
func (api *Handler) ServiceStatus(w http.ResponseWriter, r *http.Request) {
	status, err := api.usecase.GetServiceStatus()
	if err != nil {
		log.Println(err)
		ReturnErrorJSON(w, model.ErrServerError500, 500)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&status)
}

// ServiceClear godoc
// @Summary Clears Service info
// @Description Clears Service info
// @ID ServiceClear
// @Accept  json
// @Produce  json
// @Tags Service
// @Success 200 {object} nil "OK"
// @Failure 500 {object} model.Error "Internal Server Error - Request is valid but operation failed at server side"
// @Router /service/clear [post]
func (api *Handler) ServiceClear(w http.ResponseWriter, r *http.Request) {
	err := api.usecase.ServiceClear()
	if err != nil {
		log.Println(err)
		ReturnErrorJSON(w, model.ErrServerError500, 500)
		return
	}

	w.WriteHeader(200)
}
