package server

import (
	"net/http"
	"text/template"
	// "github.com/ZeineI/forum/internal/models"
)

func ErrorHandler(w http.ResponseWriter, code int) {
	tmp, err := template.ParseFiles("./static/errors.html")
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
	content := &ErrorStruct{
		ErrorNum: code,
		CodeText: http.StatusText(code),
	}
	w.WriteHeader(code)
	if err1 := tmp.Execute(w, content); err1 != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
}