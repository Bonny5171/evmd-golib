package sf

import (
	"fmt"
	"strings"

	"bitbucket.org/everymind/gforce"
)

func GetEndpoint(e string) gforce.ForceEndpoint {
	switch strings.ToLower(e) {
	case "prerelease":
		return gforce.EndpointPrerelease
	case "test", "qa", "sandbox":
		return gforce.EndpointTest
	case "mobile":
		return gforce.EndpointMobile1
	case "custom":
		return gforce.EndpointCustom
	default:
		return gforce.EndpointProduction
	}
}

func GetEndpointURL(e string) (url string, err error) {
	sfEndpoint := GetEndpoint(e)

	url, err = gforce.GetEndpointURL(sfEndpoint)
	if err != nil {
		err = fmt.Errorf("gforce.GetEndpointURL(): %w", err)
		return "", err
	}

	return
}
