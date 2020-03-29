package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
	"shorterUrl/src/env"
	"shorterUrl/src/error"
	"shorterUrl/src/middle"
)

type App struct {
	Router      *mux.Router
	Middlewares *middle.Middleware
	config      *env.Env
}

type shortenReq struct {
	URL                 string `json:"url" validate:"nonzero"`
	ExpirationInMinutes int64  `json:"expiration_in_minutes" validate:"min=0"`
}

type shortLinkResp struct {
	ShortLink string `json:"shortLink"`
}

func (a *App) Initialize(env *env.Env) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	a.config = env
	a.Router = mux.NewRouter()
	a.Middlewares = &middle.Middleware{}
	a.initializeRoute()
}

func (a *App) initializeRoute() {
	mid := alice.New(
		a.Middlewares.LoggingHandler,
		a.Middlewares.RecoverHandler)

	a.Router.Handle("/api/shorten",
		mid.ThenFunc(a.createShortLink)).Methods("POST")

	a.Router.Handle("/api/info",
		mid.ThenFunc(a.getShortLinkInfo)).Methods("GET")

	a.Router.Handle("/{shortLink:[a-zA-Z0-9]{1,11}}",
		mid.ThenFunc(a.redirect)).Methods("GET")
}

func (a *App) createShortLink(w http.ResponseWriter, r *http.Request) {
	var req shortenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseWithError(w, er.StatusError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("parse parameter failed %v", r.Body),
		})
		return
	}
	if err := validator.Validate(req); err != nil {
		responseWithError(w, er.StatusError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("validator parameter failed %v", req),
		})
		return
	}
	defer r.Body.Close()

	res, err := a.config.S.Shorten(req.URL, req.ExpirationInMinutes)
	if err != nil {
		responseWithError(w, err)
	} else {
		responseWithJSON(w, http.StatusCreated, shortLinkResp{ShortLink: res})
	}
}

func (a *App) getShortLinkInfo(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	link := vals.Get("shortLink")
	fmt.Printf("getShortLink:%s\n", link)

	add, err := a.config.S.ShortLinkInfo(link)
	if err != nil {
		responseWithError(w, err)
	} else {
		responseWithJSON(w, http.StatusOK, add)
	}
}

func (a *App) redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Printf("redirect %s\n", vars["shortLink"])
	res, err := a.config.S.UnShorten(vars["shortLink"])
	if err != nil {
		responseWithError(w, err)
	} else {
		http.Redirect(w, r, res, http.StatusTemporaryRedirect)
	}
}

func (a *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, a.Router))
}

func responseWithError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case er.Error:
		log.Printf("HTTP %d-%s", e.Status(), e)
		responseWithJSON(w, e.Status(), e.Error())
	default:
		responseWithJSON(
			w,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
	}
}

func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	resp, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}
