package googleadwords

import (
	"context"
	"fmt"
	"net/http"

	"github.com/leapforce-libraries/gads"
	bigquerytools "github.com/leapforce-libraries/go_bigquerytools"
	errortools "github.com/leapforce-libraries/go_errortools"
	types "github.com/leapforce-libraries/go_types"
	"golang.org/x/oauth2"

	go_oauth2 "github.com/leapforce-libraries/go_oauth2"
)

const (
	apiName string = "GoogleAdWords"
	//apiURL          string = "https://www.googleapis.com/calendar/v3"
	authURL         string = "https://accounts.google.com/o/oauth2/v2/auth"
	tokenURL        string = "https://oauth2.googleapis.com/token"
	tokenHTTPMethod string = http.MethodPost
	redirectURL     string = "http://localhost:8080/oauth/redirect"
)

// GoogleAdWords stores GoogleAdWords configuration
//
type GoogleAdWords struct {
	developerToken string
	oAuth2         *go_oauth2.OAuth2
}

// methods
//
func NewGoogleAdWords(developerToken string, clientID string, clientSecret string, scope string, bigQuery *bigquerytools.BigQuery, isLive bool) (*GoogleAdWords, *errortools.Error) {
	gaw := GoogleAdWords{}
	gaw.developerToken = developerToken

	config := go_oauth2.OAuth2Config{
		ApiName:         apiName,
		ClientID:        clientID,
		ClientSecret:    clientSecret,
		Scope:           scope,
		RedirectURL:     redirectURL,
		AuthURL:         authURL,
		TokenURL:        tokenURL,
		TokenHTTPMethod: tokenHTTPMethod,
	}
	gaw.oAuth2 = go_oauth2.NewOAuth(config, bigQuery, isLive)
	return &gaw, nil

}

func (gaw *GoogleAdWords) Validate() error {
	if gaw.developerToken == "" {
		return &types.ErrorString{fmt.Sprintf("%s developerToken not provided", apiName)}
	}

	return nil
}

func (gaw *GoogleAdWords) GetHttpClient() (*http.Client, *errortools.Error) {
	_, e := gaw.oAuth2.ValidateToken()
	if e != nil {
		return nil, e
	}

	return new(http.Client), nil
}

func (gaw *GoogleAdWords) GetCampaignName(customerId string, campaignId string) (string, *errortools.Error) {
	t, e := gaw.oAuth2.ValidateToken()
	if e != nil {
		return "", e
	}

	token := oauth2.Token{}
	token.AccessToken = *t.AccessToken
	token.TokenType = *t.TokenType
	token.RefreshToken = *t.RefreshToken
	token.Expiry = *t.Expiry

	authConf, _ := gads.NewCredentialsFromCode(context.TODO(), customerId, gaw.developerToken, "Leapforce", &token)

	cs := gads.NewCampaignService(&authConf.Auth)

	campaigns, _, err := cs.Get(
		gads.Selector{
			Fields: []string{
				"Id",
				"Name",
			},
			Predicates: []gads.Predicate{
				{"Id", "EQUALS", []string{campaignId}},
			},
		},
	)
	if err != nil {
		return "?", errortools.ErrorMessage(err)
	}

	if len(campaigns) > 0 {
		return campaigns[0].Name, nil
	} else {
		return "?", nil
	}
}

func (gaw *GoogleAdWords) GetCampaigns(customerId string) (*[]gads.Campaign, *errortools.Error) {
	t, e := gaw.oAuth2.ValidateToken()
	if e != nil {
		return nil, e
	}

	token := oauth2.Token{}
	token.AccessToken = *t.AccessToken
	token.TokenType = *t.TokenType
	token.RefreshToken = *t.RefreshToken
	token.Expiry = *t.Expiry

	authConf, _ := gads.NewCredentialsFromCode(context.TODO(), customerId, gaw.developerToken, "Leapforce", &token)

	cs := gads.NewCampaignService(&authConf.Auth)

	campaigns, _, err := cs.Get(
		gads.Selector{
			Fields: []string{
				"Id",
				"Name",
			},
		},
	)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}

	return &campaigns, nil
}
