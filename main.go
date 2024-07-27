package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/mux"
)

//go:embed templates/*
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

var templates = template.Must(template.ParseFS(templateFS, "templates/*.html"))

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/editor", editorHandler)
	r.HandleFunc("/run", runCodeHandler)
	r.PathPrefix("/static/").Handler(http.FileServer(http.FS(staticFS)))

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func editorHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "editor.html", nil)
}

func runCodeHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	language := r.FormValue("language")

	output, err := executeCode(code, language)
	if err != nil {
		output = fmt.Sprintf("Error: %v", err)
	}

	templates.ExecuteTemplate(w, "output.html", map[string]string{"Output": output})
}

func executeCode(code, language string) (string, error) {
	switch language {
	case "python":
		cmd := exec.Command("python", "-c", code)
		out, err := cmd.CombinedOutput()
		return string(out), err
	case "javascript":
		cmd := exec.Command("node", "-e", code)
		out, err := cmd.CombinedOutput()
		return string(out), err
	case "go":
		tempFile := fmt.Sprintf("/tmp/code_%d.go", time.Now().UnixNano())
		if err := os.WriteFile(tempFile, []byte(code), 0644); err != nil {
			return "", err
		}
		cmd := exec.Command("go", "run", tempFile)
		out, err := cmd.CombinedOutput()
		return string(out), err
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}
}
