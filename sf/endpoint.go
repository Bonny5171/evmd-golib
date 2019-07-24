package sf

import (
	"strings"

	force "bitbucket.org/everymind/gforce/lib"
	"github.com/pkg/errors"
)

func GetEndpoint(e string) force.ForceEndpoint {
	switch strings.ToLower(e) {
	case "prerelease":
		return force.EndpointPrerelease
	case "test", "qa":
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
		err = errors.Wrap(err, "force.GetEndpointURL()")
		return "", err
	}

	return
}
