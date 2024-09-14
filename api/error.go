package api

import "net/http"

// ErrorData holds the data to pass to the error template
type ErrorData struct {
    Code        int
    Message     string
    Description string
}

// Helper function to render the custom error template
func RenderErrorPage(w http.ResponseWriter, code int, message, description string) {
    w.WriteHeader(code)
    data := ErrorData{
        Code:        code,
        Message:     message,
        Description: description,
    }
    renderTemplate(w, "error.html", data)
}
