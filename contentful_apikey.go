package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/levigross/grequests"
)

type apiKeyProperties struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	AccessToken string `json:"accessToken"`
	Policies    []struct {
		Effect  string `json:"effect"`
		Actions string `json:"actions"`
	} `json:"policies"`
	PreviewAPIKey link `json:"preview_api_key"`
}

type apiKeyData struct {
	spaceEntitySys
	apiKeyProperties
}

type apiKeyCollection struct {
	collectionProperties
	Items []apiKeyData `json:"items"`
}

func apiKeyURL(baseURL, spaceID, apiKeyID string) string {
	return fmt.Sprintf("%s/spaces/%s/api_keys/%s", baseURL, spaceID, apiKeyID)
}

func createAPIKey(cmaToken, spaceID, name string) (*apiKeyData, error) {
	URL := apiKeyURL(baseURL, spaceID, "")
	res, err := grequests.Post(URL, &grequests.RequestOptions{
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", cmaToken),
			"Content-Type":  contentfulContentType,
		},
		JSON: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.StatusCode == 401 {
		return nil, errorUnauthorized
	}

	if res.StatusCode != 201 {
		return nil, errors.New(res.String())
	}

	key := &apiKeyData{}
	if err := res.JSON(key); err != nil {
		return nil, err
	}

	return key, nil
}

func updateAPIKey(
	cmaToken,
	spaceID,
	apiKeyID string,
	version int,
	newName string,
) (*apiKeyData, error) {
	URL := apiKeyURL(baseURL, spaceID, apiKeyID)
	res, err := grequests.Put(URL, &grequests.RequestOptions{
		Headers: map[string]string{
			"Authorization":        fmt.Sprintf("Bearer %s", cmaToken),
			"Content-Type":         contentfulContentType,
			"X-Contentful-Version": strconv.Itoa(version),
		},
		JSON: map[string]string{
			"name": newName,
		},
	})
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.StatusCode == 401 {
		return nil, errorUnauthorized
	}

	if res.StatusCode != 200 {
		return nil, errors.New(res.String())
	}

	key := &apiKeyData{}
	if err := res.JSON(key); err != nil {
		return nil, err
	}

	return key, nil
}

func readAPIKey(cmaToken, spaceID, apiKeyID string) (*apiKeyData, error) {
	URL := apiKeyURL(baseURL, spaceID, apiKeyID)
	res, err := grequests.Get(URL, &grequests.RequestOptions{
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", cmaToken),
		},
	})
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.StatusCode == 401 {
		return nil, errorUnauthorized
	}

	if res.StatusCode != 200 {
		return nil, errors.New(res.String())
	}

	key := &apiKeyData{}
	if err := res.JSON(key); err != nil {
		return nil, err
	}

	return key, nil
}

func deleteAPIKey(cmaToken, spaceID, apiKeyID string) error {
	URL := apiKeyURL(baseURL, spaceID, apiKeyID)
	res, err := grequests.Delete(URL, &grequests.RequestOptions{
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", cmaToken),
		},
	})
	if err != nil {
		return err
	}
	defer res.Close()

	if res.StatusCode != 204 {
		return errors.New(res.String())
	}

	return nil
}

func listAPIKeys(cmaToken, spaceID string) ([]apiKeyData, error) {
	URL := apiKeyURL(baseURL, spaceID, "")

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

	col := &apiKeyCollection{}
	if err := resp.JSON(col); err != nil {
		return nil, err
	}

	return col.Items, nil
}
