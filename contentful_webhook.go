package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/levigross/grequests"
)

type spaceEntitySys struct {
	Sys struct {
		Type      string `json:"type"`
		ID        string `json:"id"`
		Version   int    `json:"version"`
		Space     link   `json:"space"`
		CreatedAt string `json:"createdAt"`
		CreatedBy link   `json:"createdBy"`
		UpdatedAt string `json:"updatedAt"`
		UpdatedBy link   `json:"updatedBy"`
	} `json:"sys"`
}

type headerKeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type webhookProperties struct {
	Name              string   `json:"name"`
	URL               string   `json:"url"`
	Topics            []string `json:"topics"`
	HTTPBasicUsername string   `json:"httpBasicUsername"`
	// According to the CMA docs the stuff below shouldn't exist when GETing a webhook
	// TODO check
	HTTPBasicPassword string           `json:"httpBasicPassword"`
	Headers           []headerKeyValue `json:"headers"`
}

type webhookData struct {
	spaceEntitySys
	webhookProperties
}

type webhookCollection struct {
	collectionProperties
	Items []webhookData `json:"items"`
}

// func (wh *webhook) update(properties webhookProperties) error {
// updatedWh, err := updateWebhook(
// wh.client.cmaToken,
// wh.Sys.Space.Sys.ID,
// wh.Sys.Version,
// properties,
// )
// if err != nil {
// return err
// }
// wh.Sys = updatedWh.Sys
// }

func webhookPath(baseURL, spaceID, webhookID string) string {
	return fmt.Sprintf("%s/spaces/%s/webhook_definitions/%s", baseURL, spaceID, webhookID)
}

func createWebhook(
	cmaToken string,
	spaceID string,
	webhookProps webhookProperties,
) (*webhookData, error) {

	path := webhookPath(baseURL, spaceID, "")
	authHeader := fmt.Sprintf("Bearer %s", cmaToken)

	resp, err := grequests.Post(path, &grequests.RequestOptions{
		Headers: map[string]string{
			"Authorization": authHeader,
			"Content-Type":  "application/vnd.contentful.management.v1+json",
		},
		JSON: webhookProps,
	})

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 401 {
		return nil, errorUnauthorized
	}

	if resp.StatusCode == 404 {
		return nil, errorSpaceNotFound
	}

	if resp.StatusCode != 201 {
		return nil, errors.New(resp.String())
	}

	var wh webhookData
	if err := resp.JSON(&wh); err != nil {
		return nil, err
	}

	return &wh, nil
}

func updateWebhook(
	cmaToken string,
	spaceID string,
	webhookID string,
	version int,
	webhookProps webhookProperties,
) (*webhookData, error) {

	path := webhookPath(baseURL, spaceID, webhookID)
	authHeader := fmt.Sprintf("Bearer %s", cmaToken)

	resp, err := grequests.Put(path, &grequests.RequestOptions{
		Headers: map[string]string{
			"Authorization":        authHeader,
			"Content-Type":         "application/vnd.contentful.management.v1+json",
			"X-Contentful-Version": strconv.Itoa(version),
		},
		JSON: webhookProps,
	})

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 401 {
		return nil, errorUnauthorized
	}

	if resp.StatusCode == 404 {
		return nil, errorSpaceNotFound
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.String())
	}

	var wh webhookData
	if err := resp.JSON(&wh); err != nil {
		return nil, err
	}

	return &wh, nil
}

func readWebhook(cmaToken, spaceID, webhookID string) (*webhookData, error) {
	path := webhookPath(baseURL, spaceID, webhookID)
	authHeader := fmt.Sprintf("Bearer %s", cmaToken)

	res, err := grequests.Get(path, &grequests.RequestOptions{
		Headers: map[string]string{
			"Authorization": authHeader,
		},
	})

	if err != nil {
		return nil, err
	}

	if res.StatusCode == 401 {
		return nil, errorUnauthorized
	}

	// ...Sigh... The routing returns 503 if a space cannot be found, not 404.
	// It's either abiding by the broken status code, or using the /token endpoint, which
	// is more accurate but annoying. TODO: Check the /token endpoint
	if res.StatusCode == 503 {
		return nil, errorSpaceNotFound
	}

	if res.StatusCode == 404 {
		return nil, errorWebhookNotFound
	}

	if res.StatusCode != 200 {
		return nil, errors.New(res.String())
	}

	wh := &webhookData{}
	if err := res.JSON(wh); err != nil {
		return nil, err
	}

	return wh, nil
}

func deleteWebhook(cmaToken, spaceID, webhookID string) error {
	path := webhookPath(baseURL, spaceID, webhookID)
	authHeader := fmt.Sprintf("Bearer %s", cmaToken)

	resp, err := grequests.Delete(path, &grequests.RequestOptions{
		Headers: map[string]string{
			"Authorization": authHeader,
		},
	})

	if err != nil {
		return err
	}

	if resp.StatusCode == 401 {
		return errorUnauthorized
	}

	if resp.StatusCode == 404 {
		return errorSpaceNotFound
	}

	if resp.StatusCode != 204 {
		msg := fmt.Sprintf("Error: unknown status code: %d", resp.StatusCode)
		return errors.New(msg)
	}

	return nil
}

func listWebhooks(cmaToken, spaceID string) ([]webhookData, error) {
	URL := webhookPath(baseURL, spaceID, "")

	resp, err := grequests.Get(URL, &grequests.RequestOptions{
		Headers: map[string]string{
			"Authorization": "Bearer " + cmaToken,
		},
	})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 401 {
		return nil, errorUnauthorized
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.String())
	}

	webhookCol := &webhookCollection{}
	if err := resp.JSON(webhookCol); err != nil {
		return nil, err
	}

	return webhookCol.Items, nil
}
