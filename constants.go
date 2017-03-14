package main

import (
	"errors"
)

var (
	baseURL               = "https://api.contentful.com"
	contentfulContentType = "application/vnd.contentful.management.v1+json"
	// User friendly errors we return
	errorUnauthorized         = errors.New("401 Unauthorized. Is the CMA token valid?")
	errorSpaceNotFound        = errors.New("Space not found")
	errorOrganizationNotFound = errors.New("Organization not found")
	errorLocaleNotFound       = errors.New("Locale not found")
	errorWebhookNotFound      = errors.New("The webhook could not be found")
)
