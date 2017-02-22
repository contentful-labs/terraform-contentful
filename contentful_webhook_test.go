package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/stretchr/testify/assert"

	"testing"
)

type testHandler struct {
	handleFunc func(rw http.ResponseWriter, req *http.Request)
}

func (h *testHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h.handleFunc(rw, req)
	defer req.Body.Close()
}

func createExampleWebhook() webhookData {
	return webhookData{
		spaceEntitySys{},
		webhookProperties{
			Name:              "test-name",
			URL:               "http://test",
			Topics:            []string{"Entry.publish"},
			HTTPBasicUsername: "test-username",
			HTTPBasicPassword: "test-password",
			Headers: []headerKeyValue{
				headerKeyValue{
					Key:   "test-key",
					Value: "test-value",
				},
			},
		},
	}
}

func TestCreateWebhook(t *testing.T) {
	cmaToken := "test-token"
	spaceID := "test-space"

	whProps := createExampleWebhook().webhookProperties

	serveFn := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Method, "POST")
		assert.Equal(t, req.RequestURI, "/spaces/"+spaceID+"/webhook_definitions/")
		assert.Equal(t, "Bearer "+cmaToken, req.Header.Get("Authorization"))
		assert.Equal(t, "application/vnd.contentful.management.v1+json", req.Header.Get("Content-Type"))
		w.WriteHeader(201)

		var requestWh webhookData
		assert.Nil(t, json.NewDecoder(req.Body).Decode(&requestWh))
		assert.Equal(t, whProps, requestWh.webhookProperties)

		err := json.NewEncoder(w).Encode(&webhookData{})
		assert.Nil(t, err)
	}

	server := httptest.NewServer(&testHandler{handleFunc: serveFn})
	defer server.Close()
	baseURL = server.URL
	_, err := createWebhook(cmaToken, spaceID, whProps)
	assert.Nil(t, err)
}

func TestUpdateWebhook(t *testing.T) {
	cmaToken := "test-token"
	spaceID := "test-space"
	webhookID := "test-webhook"
	version := 1

	whProps := createExampleWebhook().webhookProperties

	serveFn := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Method, "PUT")
		assert.Equal(t, req.RequestURI, "/spaces/"+spaceID+"/webhook_definitions/"+webhookID)
		assert.Equal(t, "Bearer "+cmaToken, req.Header.Get("Authorization"))
		assert.Equal(t, "application/vnd.contentful.management.v1+json", req.Header.Get("Content-Type"))
		assert.Equal(t, strconv.Itoa(version), req.Header.Get("X-Contentful-Version"))
		w.WriteHeader(200)

		var requestWh webhookData
		assert.Nil(t, json.NewDecoder(req.Body).Decode(&requestWh))
		assert.Equal(t, whProps, requestWh.webhookProperties)

		err := json.NewEncoder(w).Encode(&webhookData{})
		assert.Nil(t, err)
	}

	server := httptest.NewServer(&testHandler{handleFunc: serveFn})
	defer server.Close()
	baseURL = server.URL
	_, err := updateWebhook(cmaToken, spaceID, webhookID, version, whProps)
	assert.Nil(t, err)
}

func TestReadWebhook(t *testing.T) {
	cmaToken := "test-token"
	spaceID := "test-space"
	webhookID := "test-webhook"

	wh := createExampleWebhook()

	serveFn := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Method, "GET")
		assert.Equal(t, req.RequestURI, "/spaces/"+spaceID+"/webhook_definitions/"+webhookID)
		assert.Equal(t, "Bearer "+cmaToken, req.Header.Get("Authorization"))

		err := json.NewEncoder(w).Encode(&wh)
		assert.Nil(t, err)
	}

	server := httptest.NewServer(&testHandler{handleFunc: serveFn})
	defer server.Close()
	baseURL = server.URL
	_, err := readWebhook(cmaToken, spaceID, webhookID)
	assert.Nil(t, err)
}

func TestDeleteWebhook(t *testing.T) {
	cmaToken := "test-token"
	spaceID := "test-space"
	webhookID := "test-webhook"

	serveFn := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Method, "DELETE")
		assert.Equal(t, req.RequestURI, "/spaces/"+spaceID+"/webhook_definitions/"+webhookID)
		assert.Equal(t, "Bearer "+cmaToken, req.Header.Get("Authorization"))

		w.WriteHeader(204)
	}

	server := httptest.NewServer(&testHandler{handleFunc: serveFn})
	defer server.Close()
	baseURL = server.URL
	err := deleteWebhook(cmaToken, spaceID, webhookID)
	assert.Nil(t, err)
}
