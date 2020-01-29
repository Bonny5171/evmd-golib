package sf

import (
	"fmt"
	"strings"

	force "bitbucket.org/everymind/gforce/lib"
)

func GetEndpoint(e string) force.ForceEndpoint {
	switch strings.ToLower(e) {
	case "prerelease":
		return force.EndpointPrerelease
	case "test", "qa", "sandbox":
		return force.EndpointTest
	case "mobile":
		return force.EndpointMobile1
	case "custom":
		return force.EndpointCustom
	default:
		return force.EndpointProduction
	}
}

func GetEndpointURL(e string) (url string, err error) {
	sfEndpoint := GetEndpoint(e)

	url, err = force.GetEndpointURL(sfEndpoint)
	if err != nil {
		err = fmt.Errorf("force.GetEndpointURL(): %w", err)
		return "", err
	}

	return
}
