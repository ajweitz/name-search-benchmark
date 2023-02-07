package restapp

import (
	"io"
	"log"
	"net/http"
)

func errorResponse(w http.ResponseWriter, err error, errMessage string) {
	log.Println(errMessage)
	log.Printf("Error: %v\n", err)

	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, "")
}
