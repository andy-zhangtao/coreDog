package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/andy-zhangtao/coreDog/handler"
	"github.com/andy-zhangtao/coreDog/util"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	s := r.PathPrefix("/v1").Subrouter()
	s.HandleFunc(util.RESTART_SRV, handler.RestartService).Methods(http.MethodPost)

	fmt.Printf("CoreDog[%s]准备启动\n", version())

	fmt.Println(http.ListenAndServe(":"+os.Getenv("RUNTIME_PORT"), r))
}

func version() string {
	return "v0.1"
}
