package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/andy-zhangtao/coreDog/handler"
	"github.com/andy-zhangtao/coreDog/util"
	"github.com/gorilla/mux"
)

var _VERSION_ = "unknown"

func main() {

	r := mux.NewRouter()
	s := r.PathPrefix("/v1").Subrouter()
	// s.HandleFunc(util.RESTART_SRV, handler.RestartService).Methods(http.MethodPost)
	s.HandleFunc(util.LISTALLSRV, handler.ListService).Methods(http.MethodGet)
	s.HandleFunc(util.PULLIMG, handler.PullImgService).Methods(http.MethodPut)
	s.HandleFunc(util.STARTSRV, handler.StartService).Methods(http.MethodPost)
	s.HandleFunc(util.RESTAETSRV, handler.RestartService).Methods(http.MethodPut)

	fmt.Printf("CoreDog[%s]启动\n", _VERSION_)

	fmt.Println(http.ListenAndServe(":"+os.Getenv("RUNTIME_PORT"), r))
}
