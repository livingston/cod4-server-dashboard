package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/livingston/cod4-server-dashboard/parser"
	"github.com/spf13/viper"
)

type dashboard struct {
	Title  string
	Game   map[string]string
	Server parser.Server
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("layout.html"))

	game, server, err := parser.Parse(viper.GetString("gameserver.server_location") + "serverstatus.xml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	data := dashboard{
		Title:  "Dashboard",
		Game:   game,
		Server: server,
	}

	tmpl.Execute(w, data)
}

func newRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", handler).Methods("GET")

	staticFileDirectory := http.Dir("./static/")
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

	return r
}

func main() {
	loadConfig()

	appAddress := viper.GetString("app.ip") + ":" + viper.GetString("app.port")

	r := newRouter()

	if err := http.ListenAndServe(appAddress, r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
