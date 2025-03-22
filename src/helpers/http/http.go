package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HttpClient struct {
	h http.Client
	domain string
	auth string
}

func NewClient(serverAddress string) *HttpClient {
	return &HttpClient{
		h: http.Client{},
		domain: serverAddress,
	}
}

func Get[T any](client *HttpClient, path string) (*T, error) {
	getRequest, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", client.domain, path), nil)
	if err != nil {
		return nil, err
	}

	getRequest.Header.Add("Authorization", client.auth)

	resp, err := client.h.Do(getRequest)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var data T
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func Post[T any, U any](client *HttpClient,path string, body T) (*U, error) {
	postRequest, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", client.domain, path), nil)
	if err != nil {
		return nil, err
	}

	postRequest.Header.Add("Authorization", client.auth)
	postRequest.Header.Add("Content-Type", "application/json")

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	postRequest.Body = io.NopCloser(bytes.NewReader(jsonBody))

	resp, err := client.h.Do(postRequest)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var result U
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func Delete[T any](client *HttpClient, path string) (*T, error) {
	getRequest, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", client.domain, path), nil)
	if err != nil {
		return nil, err
	}

	getRequest.Header.Add("Authorization", client.auth)

	resp, err := client.h.Do(getRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return nil, nil
}

func Patch[T any, U any](client *HttpClient, path string, body T) (*U, error) {
	postRequest, err := http.NewRequest("PATCH", fmt.Sprintf("%s/%s", client.domain, path), nil)
	if err != nil {
		return nil, err
	}

	postRequest.Header.Add("Authorization", client.auth)
	postRequest.Header.Add("Content-Type", "application/json")

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	postRequest.Body = io.NopCloser(bytes.NewReader(jsonBody))

	resp, err := client.h.Do(postRequest)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var result U
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}