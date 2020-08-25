package googleadwords

import (
	"fmt"
	"net/http"
	"strings"

	bigquerytools "github.com/Leapforce-nl/go_bigquerytools"
	types "github.com/Leapforce-nl/go_types"

	googleoauth2 "github.com/Leapforce-nl/go_googleoauth2"
)

const apiName string = "GoogleAdWords"

// GoogleAdWords stores GoogleAdWords configuration
//
type GoogleAdWords struct {
	SiteURL string
	BaseURL string
	oAuth2  *googleoauth2.GoogleOAuth2
}

// methods
//
func (gaw *GoogleAdWords) InitOAuth2(clientID string, clientSecret string, bigQuery *bigquerytools.BigQuery, isLive bool) error {
	_oAuth2 := new(googleoauth2.GoogleOAuth2)
	_oAuth2.ApiName = apiName
	_oAuth2.ClientID = clientID
	_oAuth2.ClientSecret = clientSecret
	_oAuth2.BigQuery = bigQuery
	_oAuth2.IsLive = isLive

	gaw.oAuth2 = _oAuth2

	return nil
}

func (gaw *GoogleAdWords) Validate() error {
	if gaw.BaseURL == "" {
		return &types.ErrorString{fmt.Sprintf("%s BaseURL not provided", apiName)}
	}
	if gaw.SiteURL == "" {
		return &types.ErrorString{fmt.Sprintf("%s SiteURL not provided", apiName)}
	}

	if !strings.HasSuffix(gaw.BaseURL, "/") {
		gaw.BaseURL = gaw.BaseURL + "/"
	}

	if !strings.HasSuffix(gaw.SiteURL, "/") {
		gaw.SiteURL = gaw.SiteURL + "/"
	}

	return nil
}

func (gaw *GoogleAdWords) GetHttpClient() (*http.Client, error) {

	err := gaw.oAuth2.ValidateToken()
	if err != nil {
		return nil, err
	}

	return new(http.Client), nil
}
