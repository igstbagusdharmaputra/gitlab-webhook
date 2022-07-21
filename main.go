package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func MyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiKey := os.Getenv("API_KEY")
		header := r.Header.Get("X-Gitlab-Token")

		if header != apiKey {
			APIResponse(w, "Status Unauthorized", http.StatusUnauthorized, rError, "")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func sendMessage(data Gitlab) error {
	request_url := getSendMessageURL()
	client := &http.Client{}

	if data.ObjectKind == "pipeline" {
		if data.ObjectAttributes.Status == "pending" {
			rawMsg := fmt.Sprintf("CI/CD Run Process with User Name *%s* Project Namespace *%s* Project Name *%s*", data.User.Name, data.Project.Namespace, data.Project.Name)
			reqBody := sendMessageReqBody{
				Text: rawMsg,
			}

			reqBytes, err := json.Marshal(&reqBody)
			if err != nil {
				log.Printf("Json unmarshal error, %v", err)
			}
			req, _ := http.NewRequest("POST", request_url, bytes.NewBuffer(reqBytes))
			req.Header.Set("Content-Type", "application/json")
			res, err := client.Do(req)

			if err != nil {
				log.Printf("Message can't send, %v", err)
			} else {
				defer res.Body.Close()
			}
		} else {
			var rawMsg string
			for _, v := range data.Builds {
				if v.Status == "success" {
					rawMsg = fmt.Sprintf("✅ The job *%s* for project *%s/%s* was run by *%s*. Job completed successfully", v.User.Name, data.Project.Namespace, data.Project.Name, v.Name)
				} else if v.Status == "failed" {
					rawMsg = fmt.Sprintf("❌ The job *%s* for project *%s/%s* was run by *%s*. Job Failed", v.User.Name, data.Project.Namespace, data.Project.Name, v.Name)
				} else if v.Status == "skipped" {
					rawMsg = fmt.Sprintf("⛔ The job *%s* for project *%s/%s* was run by *%s*. Job Skipped", v.User.Name, data.Project.Namespace, data.Project.Name, v.Name)
				} else {
					rawMsg = ""
				}

				reqBody := sendMessageReqBody{
					Text: rawMsg,
				}

				reqBytes, err := json.Marshal(&reqBody)
				if err != nil {
					log.Printf("Json unmarshal error, %v", err)
				}
				req, _ := http.NewRequest("POST", request_url, bytes.NewBuffer(reqBytes))
				req.Header.Set("Content-Type", "application/json")
				res, err := client.Do(req)

				if err != nil {
					log.Printf("Message can't send, %v", err)
				} else {
					log.Println(res.Status)
					defer res.Body.Close()
				}
			}
		}
	}
	return nil

}
func webhook(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	data := Gitlab{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		APIResponse(w, "Status BadRequest", http.StatusBadRequest, rError, "")
		return
	}
	sendMessage(data)
	APIResponse(w, "Status OK", http.StatusOK, rSuccess, data)
}
func main() {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	log.Printf("Running on port: %d", port)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.With(MyMiddleware).Post("/webhook", webhook)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
