package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	// "github.com/parnurzeal/gorequest"
	"github.com/stretchr/testify/assert"
)

func TestCreateApikey(t *testing.T) {
	cmaToken := "test-token"
	spaceID := "test-space"
	apikeyName := "test-apikey"

	serveFn := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Method, "POST")
		assert.Equal(t, req.RequestURI, "/spaces/"+spaceID+"/api_keys/")
		assert.Equal(t, "Bearer "+cmaToken, req.Header.Get("Authorization"))
		assert.Equal(t, "application/vnd.contentful.management.v1+json", req.Header.Get("Content-Type"))

		var payload map[string]string
		assert.Nil(t, json.NewDecoder(req.Body).Decode(&payload))
		assert.Equal(t, apikeyName, payload["name"])

		w.WriteHeader(201)
		err := json.NewEncoder(w).Encode(&apiKeyData{})
		assert.Nil(t, err)
	}

	server := httptest.NewServer(&testHandler{handleFunc: serveFn})
	defer server.Close()
	baseURL = server.URL

	_, err := createAPIKey(cmaToken, spaceID, apikeyName)
	assert.Nil(t, err)
}

func TestUpdateApikey(t *testing.T) {
	cmaToken := "test-token"
	spaceID := "test-space"
	apiKeyID := "test-apikey"
	newAPIKeyName := "test-apikey"
	version := 1

	serveFn := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Method, "PUT")
		assert.Equal(t, req.RequestURI, "/spaces/"+spaceID+"/api_keys/"+apiKeyID)
		assert.Equal(t, "Bearer "+cmaToken, req.Header.Get("Authorization"))
		assert.Equal(t, "application/vnd.contentful.management.v1+json", req.Header.Get("Content-Type"))
		assert.Equal(t, strconv.Itoa(version), req.Header.Get("X-Contentful-Version"))

		var payload map[string]string
		assert.Nil(t, json.NewDecoder(req.Body).Decode(&payload))
		assert.Equal(t, newAPIKeyName, payload["name"])

		w.WriteHeader(200)
		err := json.NewEncoder(w).Encode(&apiKeyData{})
		assert.Nil(t, err)
	}

	server := httptest.NewServer(&testHandler{handleFunc: serveFn})
	defer server.Close()
	baseURL = server.URL

	_, err := updateAPIKey(cmaToken, spaceID, apiKeyID, version, newAPIKeyName)
	assert.Nil(t, err)
}

func TestReadApikey(t *testing.T) {
	cmaToken := "test-token"
	spaceID := "test-space"
	apiKeyID := "test-id"

	serveFn := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Method, "GET")
		assert.Equal(t, req.RequestURI, "/spaces/"+spaceID+"/api_keys/"+apiKeyID)
		assert.Equal(t, "Bearer "+cmaToken, req.Header.Get("Authorization"))

		w.WriteHeader(200)
		err := json.NewEncoder(w).Encode(&apiKeyData{})
		assert.Nil(t, err)
	}

	server := httptest.NewServer(&testHandler{handleFunc: serveFn})
	defer server.Close()
	baseURL = server.URL

	_, err := readAPIKey(cmaToken, spaceID, apiKeyID)
	assert.Nil(t, err)
}

func TestDeleteApikey(t *testing.T) {
	cmaToken := "test-token"
	spaceID := "test-space"
	apiKeyID := "test-id"

	serveFn := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Method, "DELETE")
		assert.Equal(t, req.RequestURI, "/spaces/"+spaceID+"/api_keys/"+apiKeyID)
		assert.Equal(t, "Bearer "+cmaToken, req.Header.Get("Authorization"))

		w.WriteHeader(204)
	}

	server := httptest.NewServer(&testHandler{handleFunc: serveFn})
	defer server.Close()
	baseURL = server.URL

	err := deleteAPIKey(cmaToken, spaceID, apiKeyID)
	assert.Nil(t, err)

}
