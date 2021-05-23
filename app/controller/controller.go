package controller

import (
	"log"
	"net/http"

	"accounts/util"
)

type Controller struct {
	service accountService
}

func New(service accountService) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) FilterAccounts(w http.ResponseWriter, r *http.Request) {
	body, err := c.service.FilterAccounts(r.Context(), r.URL.Query())
	if err != nil {
		util.WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	util.WriteSuccessResponse(w, body, http.StatusOK)
}

func (c *Controller) GroupAccounts(w http.ResponseWriter, r *http.Request) {
	log.Println("group accounts")
}

func (c *Controller) GetRecommends(w http.ResponseWriter, r *http.Request) {
	log.Println("get recommends")
}

func (c *Controller) GetSuggestions(w http.ResponseWriter, r *http.Request) {
	log.Println("get suggestions")
}

func (c *Controller) CreateAccount(w http.ResponseWriter, r *http.Request) {
	body, err := util.ReadRequestBody(r)
	if err != nil {
		util.WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err = c.service.AddAccount(r.Context(), body); err != nil {
		util.WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	util.WriteSuccessResponse(w, []byte("{}"), http.StatusCreated)
}

func (c *Controller) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	id := util.ReadURLParam(r, "id")
	if id == "" {
		util.WriteErrorResponse(w, nil, http.StatusBadRequest)
	}

	body, err := util.ReadRequestBody(r)
	if err != nil {
		util.WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err = c.service.UpdateAccount(r.Context(), body); err != nil {
		util.WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	util.WriteSuccessResponse(w, []byte("{}"), http.StatusAccepted)
}

func (c *Controller) AddLikes(w http.ResponseWriter, r *http.Request) {
	body, err := util.ReadRequestBody(r)
	if err != nil {
		util.WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err = c.service.AddLikes(r.Context(), body); err != nil {
		util.WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	util.WriteSuccessResponse(w, []byte("{}"), http.StatusAccepted)
}
