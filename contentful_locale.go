package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/levigross/grequests"
)

type localeProperties struct {
	Name         string `json:"name"`
	Code         string `json:"code"`
	FallbackCode string `json:"fallbackCode"`
	Optional     bool   `json:"optional"`
}

type locale struct {
	spaceEntitySys
	localeProperties
}

func localePath(baseURL, spaceID, localeID string) string {
	return fmt.Sprintf("%s/spaces/%s/locales/%s", baseURL, spaceID, localeID)
}

func readLocale(cmaToken, spaceID, localeID string) (*locale, error) {
	path := localePath(baseURL, spaceID, localeID)
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
		return nil, errorLocaleNotFound
	}

	if res.StatusCode != 200 {
		return nil, errors.New(res.String())
	}

	var json locale
	if err := res.JSON(&json); err != nil {
		return nil, err
	}

	return &json, nil
}

func createLocale(
	cmaToken string,
	spaceID string,
	localeProps localeProperties,
) (*locale, error) {

	path := localePath(baseURL, spaceID, "")
	authHeader := fmt.Sprintf("Bearer %s", cmaToken)

	res, err := grequests.Post(path, &grequests.RequestOptions{
		Headers: map[string]string{
			"Authorization": authHeader,
			"Content-Type":  "application/vnd.contentful.management.v1+json",
		},
		JSON: localeProps,
	})

	if err != nil {
		return nil, err
	}

	if res.StatusCode == 401 {
		return nil, errorUnauthorized
	}

	if res.StatusCode == 404 {
		return nil, errorSpaceNotFound
	}

	if res.StatusCode != 200 {

		log.Println("path", path)
		log.Println("status code", res.StatusCode)
		log.Printf("dickhead %+v\n", localeProps)
		b, _ := json.Marshal(localeProps)
		log.Println(string(b))
		log.Println(path)

		return nil, errors.New(res.String())
	}

	var json locale
	if err := res.JSON(&json); err != nil {
		return nil, err
	}

	return &json, nil
}

func updateLocale(
	cmaToken string,
	spaceID string,
	localeID string,
	localeProps localeProperties,
) (*locale, error) {

	path := localePath(baseURL, spaceID, "")
	authHeader := fmt.Sprintf("Bearer %s", cmaToken)

	res, err := grequests.Put(path, &grequests.RequestOptions{
		Headers: map[string]string{
			"Authorization": authHeader,
			"Content-Type":  "application/vnd.contentful.management.v1+json",
		},
		JSON: localeProps,
	})

	if err != nil {
		return nil, err
	}

	if res.StatusCode == 401 {
		return nil, errorUnauthorized
	}

	if res.StatusCode == 404 {
		return nil, errorSpaceNotFound
	}

	if res.StatusCode != 200 {

		log.Println("path", path)
		log.Println("status code", res.StatusCode)
		log.Printf("dickhead %+v\n", localeProps)
		b, _ := json.Marshal(localeProps)
		log.Println(string(b))
		log.Println(path)

		return nil, errors.New(res.String())
	}

	var json locale
	if err := res.JSON(&json); err != nil {
		return nil, err
	}

	return &json, nil
}

func deleteLocale(cmaToken, spaceID, localeID string) error {
	path := localePath(baseURL, spaceID, localeID)
	authHeader := fmt.Sprintf("Bearer %s", cmaToken)

	res, err := grequests.Delete(path, &grequests.RequestOptions{
		Headers: map[string]string{
			"Authorization": authHeader,
		},
	})

	if err != nil {
		return err
	}

	if res.StatusCode == 401 {
		return errorUnauthorized
	}

	if res.StatusCode == 404 {
		return errorSpaceNotFound
	}

	if res.StatusCode != 204 {
		msg := fmt.Sprintf("Error: unknown status code: %d", res.StatusCode)
		return errors.New(msg)
	}

	return nil
}
