package main

import (
	"log"
	"net/http"
)

// функция-обработчик
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}

func main() {
	// роутер
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	log.Println("Starting server on :4000\nhttp://localhost:4000/")
	// запускаем новый веб-сервер через функциюю http.ListenAndServe
	// передаем два аргумента: TCP-адрес сети для прослушивания, роутер
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
