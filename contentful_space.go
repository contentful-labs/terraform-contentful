package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/levigross/grequests"
)

type collectionProperties struct {
	Total int `json:"total"`
	Limit int `json:"limit"`
	Skip  int `json:"skip"`
	Sys   struct {
		Type string `json:"type"`
	} `json:"sys"`
}

type link struct {
	Sys struct {
		Type     string `json:"type"`
		LinkType string `json:"linkType"`
		ID       string `json:"id"`
	} `json:"sys"`
}

type spaceSys struct {
	Type      string `json:"type"`
	ID        string `json:"id"`
	Version   int    `json:"version"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	CreatedBy link   `json:"createdBy"`
	UpdatedBy link   `json:"updatedBy"`
}

type spaceData struct {
	Sys  spaceSys `json:"sys"`
	Name string   `json:"name"`
}

type spaceCollection struct {
	collectionProperties
	Items []spaceData `json:"items"`
}

func createSpace(
	cmaToken string,
	organizationID string,
	spaceName string,
	defaultLocale string,
) (*spaceData, error) {
	URL := baseURL + "/spaces/"
	res, err := grequests.Post(URL, &grequests.RequestOptions{
		Headers: map[string]string{
			"Authorization":             fmt.Sprintf("Bearer %s", cmaToken),
			"Content-Type":              contentfulContentType,
			"X-Contentful-Organization": organizationID,
		},
		JSON: map[string]string{
			"name":          spaceName,
			"defaultLocale": defaultLocale,
		},
	})
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.StatusCode == 401 {
		return nil, errorUnauthorized
	}

	if res.StatusCode == 404 {
		return nil, errorOrganizationNotFound
	}

	if res.StatusCode != 201 {
		return nil, errors.New(res.String())
	}

	s := &spaceData{}
	if err := res.JSON(s); err != nil {
		return nil, err
	}

	return s, nil
}

func updateSpace(
	cmaToken string,
	spaceID string,
	spaceVersion int,
	newSpaceName string,
) (*spaceData, error) {
	URL := baseURL + "/spaces/" + spaceID
	res, err := grequests.Put(URL, &grequests.RequestOptions{
		Headers: map[string]string{
			"Authorization":        fmt.Sprintf("Bearer %s", cmaToken),
			"Content-Type":         contentfulContentType,
			"X-Contentful-Version": strconv.Itoa(spaceVersion),
		},
		JSON: map[string]string{
			"name": newSpaceName,
		},
	})
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.StatusCode == 401 {
		return nil, errorUnauthorized
	}

	// ...Sigh...
	if res.StatusCode == 503 {
		return nil, errorSpaceNotFound
	}

	if res.StatusCode != 200 {
		return nil, errors.New(res.String())
	}

	s := &spaceData{}
	if err := res.JSON(s); err != nil {
		return nil, err
	}

	return s, nil
}

func readSpace(cmaToken, spaceID string) (*spaceData, error) {
	URL := baseURL + "/spaces/" + spaceID
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

	// ...Sigh... The routing returns 503 if a space cannot be found, not 404.
	// It's either abiding by the broken status code, or using the /token endpoint, which
	// is more accurate but annoying. TODO: Check the /token endpoint
	if res.StatusCode == 503 {
		return nil, errorSpaceNotFound
	}

	if res.StatusCode == 404 {
		return nil, errorSpaceNotFound
	}

	if res.StatusCode != 200 {
		return nil, errors.New(res.String())
	}

	s := &spaceData{}
	if err := res.JSON(s); err != nil {
		return nil, err
	}

	return s, nil
}

func deleteSpace(cmaToken, spaceID string) error {
	URL := baseURL + "/spaces/" + spaceID
	res, err := grequests.Delete(URL, &grequests.RequestOptions{
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", cmaToken),
		},
	})
	if err != nil {
		return err
	}
	defer res.Close()

	if res.StatusCode == 401 {
		return errorUnauthorized
	}

	if res.StatusCode == 404 {
		return errorSpaceNotFound
	}

	if res.StatusCode != 204 {
		return errors.New(res.String())
	}

	return nil
}

func listSpaces(cmaToken string) ([]spaceData, error) {
	URL := baseURL + "/spaces/"

	res, err := grequests.Get(URL, &grequests.RequestOptions{
		Headers: map[string]string{
			"Authorization": "Bearer " + cmaToken,
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

	spaceCol := &spaceCollection{}
	if err := res.JSON(spaceCol); err != nil {
		return nil, err
	}

	return spaceCol.Items, nil
}

func spaceExists(
	cmaToken string,
	spaceName string,
) (bool, error) {
	URL := baseURL + "/spaces/" + spaceName
	res, err := grequests.Get(URL, &grequests.RequestOptions{
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", cmaToken),
			"Content-Type":  contentfulContentType,
		},
	})
	if err != nil {
		return false, err
	}
	defer res.Close()

	if res.StatusCode == 200 {
		return true, nil
	}

	if res.StatusCode == 401 {
		return false, errorUnauthorized
	}

	if res.StatusCode != 200 && res.StatusCode != 404 {
		return false, errors.New(res.String())
	}

	return false, nil
}
