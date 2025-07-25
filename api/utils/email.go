package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gold-savings/internal/config"
	"io"
	"math/rand"
	"net/http"
	"text/template"
	"time"
)

type Plunk struct {
	HttpClient *http.Client
	Config     *config.Config
}

type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// Function to generate 4-digit OTP
func GenerateOTP() string {
	rand.NewSource(time.Now().UnixNano()) // Seed the random number generator with current time

	otp := rand.Intn(10000)         // Generate a random number between 0 to 9999
	return fmt.Sprintf("%04d", otp) // Format the OTP to always be 4 digits
}

// RenderEmailTemplate parses and executes an HTML template with the provided data.
func RenderEmailTemplate(templatePath string, data any) (string, error) {
	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	var tplBody bytes.Buffer
	if err := tpl.Execute(&tplBody, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return tplBody.String(), nil
}

func (s *Plunk) makeRequest(method, endpoint string, body any) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, s.Config.PlunkBaseUrl+endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.Config.PlunkSecretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New(string(respBody))
	}

	return respBody, nil
}

func (s *Plunk) SendEmail(to, subject, body string) error {
	email := EmailRequest{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	_, err := s.makeRequest("POST", "/send", email)
	return err
}
