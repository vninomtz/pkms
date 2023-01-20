package api

import (
	"encoding/json"
	"net/http"

	"github.com/vninomtz/swe-notes/model"
)


type noteHandler struct {
  notes []*model.Note
}


func NewNoteHandler(ns []*model.Note) *noteHandler  {
  return &noteHandler{
    notes: ns,
  }
}

type response struct {
  Body []*model.Note `json:"body"`
  Records int `json:"records"`
}

func (h *noteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  res := &response{
    Body: h.notes,
    Records: len(h.notes),
  }
  err := json.NewEncoder(w).Encode(res)
  if err != nil {
    http.Error(w, err.Error(), 500)
  }
}
