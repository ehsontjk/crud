package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"github.com/ehsontjk/crud/pkg/customers"
	"github.com/gorilla/mux"
)
type Server struct {
	mux *mux.Router
	customerSvc *customers.Service
}


func NewServer(m *mux.Router, cSvc *customers.Service) *Server {
	return &Server{mux: m, customerSvc: cSvc}
}

func (s *Server)ServeHTTP(w http.ResponseWriter, r *http.Request){
	s.mux.ServeHTTP(w,r)
}


func (s *Server) Init() {
	s.mux.HandleFunc("/customers/{id}", s.handleGetCustomerByID).Methods("GET")
	s.mux.HandleFunc("/customers", s.handleGetAllCustomers).Methods("GET")
	s.mux.HandleFunc("/customers.getAllActive", s.handleGetAllActiveCustomers)
	s.mux.HandleFunc("/customers.blockById", s.handleBlockByID)
	s.mux.HandleFunc("/customers.unblockById", s.handleUnBlockByID)
	s.mux.HandleFunc("/customers/{id}", s.handleDelete).Methods("DELETE")
	s.mux.HandleFunc("/customers", s.handleSave).Methods("POST")
}


func (s *Server) handleGetAllCustomers(w http.ResponseWriter, r *http.Request) {

	items, err :=s.customerSvc.All(r.Context())
	if err != nil{

	errorWriter(w, http.StatusInternalServerError, err)
	return
	}
	

	respondJSON(w, items)
}


func (s *Server) handleGetAllActiveCustomers(w http.ResponseWriter, r *http.Request) {

	items, err :=s.customerSvc.AllActive(r.Context())
	if err != nil{
	
	errorWriter(w, http.StatusInternalServerError, err)
	return
	}
	

	respondJSON(w, items)
}

func (s *Server) handleGetCustomerByID(w http.ResponseWriter, r *http.Request) {
	
	idP:= mux.Vars(r)["id"]

	
	id, err := strconv.ParseInt(idP, 10, 64)
	
	if err != nil {
	
	errorWriter(w, http.StatusBadRequest, err)
	return
	}

	
	item, err := s.customerSvc.ByID(r.Context(), id)

	if errors.Is(err, customers.ErrNotFound) {
	
	errorWriter(w, http.StatusNotFound, err)
	return
	}

	
	if err != nil {
	
	errorWriter(w, http.StatusInternalServerError, err)
	return
	}

	respondJSON(w, item)
}

func (s *Server) handleBlockByID(w http.ResponseWriter, r *http.Request) {
	idP := r.URL.Query().Get("id")
id, err := strconv.ParseInt(idP, 10, 64)
	if err != nil {
	errorWriter(w, http.StatusBadRequest, err)
	return
	}
item, err := s.customerSvc.ChangeActive(r.Context(), id, false)
if errors.Is(err, customers.ErrNotFound) {
	errorWriter(w, http.StatusNotFound, err)
	return
	}
if err != nil {
	
	errorWriter(w, http.StatusInternalServerError, err)
	return
	}
	respondJSON(w, item)
}

func (s *Server) handleUnBlockByID(w http.ResponseWriter, r *http.Request) {
	
	idP := r.URL.Query().Get("id")

	
	id, err := strconv.ParseInt(idP, 10, 64)
	
	if err != nil {
	
	errorWriter(w, http.StatusBadRequest, err)
	
	}

	
	item, err := s.customerSvc.ChangeActive(r.Context(), id, true)

	if errors.Is(err, customers.ErrNotFound) {
	
	errorWriter(w, http.StatusNotFound, err)
	return
	}

	
	if err != nil {
	
	errorWriter(w, http.StatusInternalServerError, err)
	return
	}
	
	respondJSON(w, item)
}

func (s *Server) handleSave(w http.ResponseWriter, r *http.Request) {
	
	var item *customers.Customer

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	
	customer, err := s.customerSvc.Save(r.Context(), item)
      if err != nil {
		
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	
	respondJSON(w, customer)
}



func errorWriter(w http.ResponseWriter, httpSts int, err error) {
	
	log.Print(err)
	
	http.Error(w, http.StatusText(httpSts), httpSts)
}

func respondJSON(w http.ResponseWriter, iData interface{}) {

	
	data, err := json.Marshal(iData)

	
	if err != nil {
		
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(data)
	
	if err != nil {
		
		log.Print(err)
	}
}
func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	
	idP := mux.Vars(r)["id"]

	
	id, err := strconv.ParseInt(idP, 10, 64)
	
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	
	item, err := s.customerSvc.Delete(r.Context(), id)
	
	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(w, http.StatusNotFound, err)
		return
	}

	
	if err != nil {
	
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	
	respondJSON(w, item)
}
