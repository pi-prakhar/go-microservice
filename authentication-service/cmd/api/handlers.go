package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go-microservice-authentication/data"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user against the database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// log authentication
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) Register(w http.ResponseWriter, r *http.Request) {

	var userPayload data.User

	err := app.readJSON(w, r, &userPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(userPayload.Email)

	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	if user != nil {
		app.errorJSON(w, errors.New("user already present"), http.StatusBadRequest)
		return
	}

	status, err := app.Models.User.Insert(userPayload)

	if status == 0 {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Successfully registered user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusOK, payload)

}
func (app *Config) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET
	if r.Method != http.MethodGet {
		app.errorJSON(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	// Call the service function to get all users
	users, err := app.Models.User.GetAll()
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Prepare the response payload
	payload := jsonResponse{
		Error:   false,
		Message: "Successfully retrieved all users",
		Data:    users,
	}

	// Write the JSON response
	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
