package controller

import (
	"log"
	"net/http"
)

type Controller struct {
}

func (c *Controller) FilterAccounts(w http.ResponseWriter, r *http.Request) {
	log.Println("filter accounts")
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
	log.Println("create account")
}

func (c *Controller) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	log.Println("update account")
}

func (c *Controller) AddLikes(w http.ResponseWriter, r *http.Request) {
	log.Println("add likes")
}
