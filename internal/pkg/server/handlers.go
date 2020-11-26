package server

import (
	"Mongo/internal/pkg/auth"
	"Mongo/internal/pkg/constErr"
	"Mongo/internal/pkg/data"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var acc data.Account

	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		http.Error(w, constErr.ErrorUnmarshal.Error(), 400)
		return
	}
	token, err := s.service.Login(&acc)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}
	w.Header().Set("Token", token)
	w.WriteHeader(201)
	w.Write([]byte("See in tab Headers."))

}

func (s *Server) SignUp(w http.ResponseWriter, r *http.Request) {
	var acc data.Account

	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		http.Error(w, constErr.ErrorUnmarshal.Error(), 400)
		return
	}

	token, err := s.service.SignUp(&acc)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	w.Header().Set("Token", token)
	fmt.Fprint(w, "See in headers refresh Token.")
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	baseAcc, err := s.service.GetStorage().GetAccounts()
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	acc, err := auth.AccountIdentification(r, baseAcc)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	err = s.service.Logout(acc)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Fprint(w, "You Logout! Bye-bye...")
}

func (s *Server) Get(w http.ResponseWriter, r *http.Request) {
	//----------------------------------------------------------------
	// token validation area
	baseAcc, err := s.service.GetStorage().GetAccounts()
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	_, err = auth.AccountIdentification(r, baseAcc)
	if err != nil {
		http.Error(w, err.Error(), 403)
		return
	}
	//-----------------------------------------------------------------
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)["id"]
	id, err := strconv.Atoi(vars)
	if err != nil && id < 0 {
		http.Error(w, constErr.ErrorConvertToInteger.Error(), 500)
		return
	}

	ad, err := s.service.Get(uint(id))
	if err != nil {
		http.Error(w, constErr.NotFoundAd.Error(), 404)
		return
	}

	js, err := json.Marshal(ad)
	if err != nil {
		http.Error(w, constErr.ErrorMarshal.Error(), 400)
		return
	}
	fmt.Fprint(w, string(js))
}

func (s *Server) GetAll(w http.ResponseWriter, r *http.Request) {
	//----------------------------------------------------------------
	// token validation area
	baseAcc, err := s.service.GetStorage().GetAccounts()
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}
	_, err = auth.AccountIdentification(r, baseAcc)
	if err != nil {
		http.Error(w, err.Error(), 403)
		return
	}
	//-----------------------------------------------------------------
	w.Header().Set("Content-Type", "application/json")
	base, err := s.service.GetAll()
	if base == nil && err != nil {
		if err == constErr.AdBaseIsEmpty {
			http.Error(w, err.Error(), http.StatusNotFound) // TODO
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	js, err := json.Marshal(base)
	if err != nil {
		http.Error(w, constErr.ErrorMarshal.Error(), 400) // TODO
		return
	}
	fmt.Fprint(w, string(js))
}

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	//----------------------------------------------------------------
	// token validation area
	baseAcc, err := s.service.GetStorage().GetAccounts()
	if err != nil {
		log.Println(1)
		http.Error(w, err.Error(), 401)
		return
	}
	_, err = auth.AccountIdentification(r, baseAcc)
	if err != nil {
		log.Println(2)
		http.Error(w, err.Error(), 403)
		return
	}
	//-----------------------------------------------------------------
	var ad data.Ad
	if err := json.NewDecoder(r.Body).Decode(&ad); err != nil {
		http.Error(w, "Error JSON!", 400)
		return
	}
	err = s.service.Add(&ad)
	if err != nil {
		http.Error(w, "Error to add.", 400)
		return
	}
	w.WriteHeader(201)
	fmt.Fprintf(w, "Create new Ad")
}

func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	//----------------------------------------------------------------
	// token and account validation area
	baseAcc, err := s.service.GetStorage().GetAccounts()
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}
	_, err = auth.AccountIdentification(r, baseAcc)
	if err != nil {
		http.Error(w, err.Error(), 403)
		return
	}
	//-----------------------------------------------------------------
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)["id"]
	id, err := strconv.Atoi(vars)
	if err != nil && id < 0 {
		http.Error(w, constErr.ErrorConvertToInteger.Error(), http.StatusBadRequest) // TODO
		return
	}
	if err := s.service.Delete(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // TODO
		return
	}
	fmt.Fprint(w, "Delete")
}

func (s *Server) Update(w http.ResponseWriter, r *http.Request) {
	//-----------------------------------------------------------------
	// token validation area
	baseAcc, err := s.service.GetStorage().GetAccounts()
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}
	_, err = auth.AccountIdentification(r, baseAcc)
	if err != nil {
		http.Error(w, err.Error(), 403)
		return
	}
	//-----------------------------------------------------------------
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)["id"]
	id, err := strconv.Atoi(vars)
	if err != nil && id < 0 {
		http.Error(w, err.Error(), http.StatusBadRequest) // TODO
		return
	}

	var ad data.Ad
	if err := json.NewDecoder(r.Body).Decode(&ad); err != nil {
		http.Error(w, constErr.ErrorUnmarshal.Error(), 400)
		return
	}

	if err = s.service.Update(uint(id), &ad); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Fprint(w, "Update")
}
