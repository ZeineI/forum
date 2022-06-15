package server

import (
	"html/template"
	"net/http"

	"github.com/ZeineI/forum/internal/models"
)

func ErrorHandler(w http.ResponseWriter, code int) {
	tmp, err := template.ParseFiles("./frontend/html/errors.html")
	if err != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
	content := &models.ErrorStruct{
		ErrorNum: code,
		CodeText: http.StatusText(code),
	}
	w.WriteHeader(code)
	if err1 := tmp.Execute(w, content); err1 != nil {
		http.Error(w, "500: Internal Server Error", http.StatusInternalServerError)
		return
	}
}
