package googleadwords

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Leapforce-nl/gads"
	bigquerytools "github.com/Leapforce-nl/go_bigquerytools"
	types "github.com/Leapforce-nl/go_types"
	"golang.org/x/oauth2"

	googleoauth2 "github.com/Leapforce-nl/go_googleoauth2"
)

const apiName string = "GoogleAdWords"

// GoogleAdWords stores GoogleAdWords configuration
//
type GoogleAdWords struct {
	DeveloperToken string
	oAuth2         *googleoauth2.GoogleOAuth2
}

// methods
//
func (gaw *GoogleAdWords) InitOAuth2(clientID string, clientSecret string, scopes []string, bigQuery *bigquerytools.BigQuery, isLive bool) error {
	_oAuth2 := new(googleoauth2.GoogleOAuth2)
	_oAuth2.ApiName = apiName
	_oAuth2.ClientID = clientID
	_oAuth2.ClientSecret = clientSecret
	_oAuth2.Scopes = scopes
	_oAuth2.BigQuery = bigQuery
	_oAuth2.IsLive = isLive

	gaw.oAuth2 = _oAuth2

	return nil
}

func (gaw *GoogleAdWords) Validate() error {
	if gaw.DeveloperToken == "" {
		return &types.ErrorString{fmt.Sprintf("%s DeveloperToken not provided", apiName)}
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

func (gaw *GoogleAdWords) GetCampaignName(customerId string, campaignId string) (string, error) {
	err := gaw.oAuth2.ValidateToken()
	if err != nil {
		return "", err
	}

	token := oauth2.Token{}
	token.AccessToken = gaw.oAuth2.Token.AccessToken
	token.TokenType = gaw.oAuth2.Token.TokenType
	token.RefreshToken = gaw.oAuth2.Token.RefreshToken
	token.Expiry = gaw.oAuth2.Token.Expiry

	authConf, _ := gads.NewCredentialsFromCode(context.TODO(), customerId, gaw.DeveloperToken, "Leapforce", &token)

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
		return "?", err
	}

	if len(campaigns) > 0 {
		return campaigns[0].Name, nil
	} else {
		return "?", nil
	}
}

func (gaw *GoogleAdWords) GetCampaigns(customerId string) (*[]gads.Campaign, error) {
	err := gaw.oAuth2.ValidateToken()
	if err != nil {
		return nil, err
	}

	token := oauth2.Token{}
	token.AccessToken = gaw.oAuth2.Token.AccessToken
	token.TokenType = gaw.oAuth2.Token.TokenType
	token.RefreshToken = gaw.oAuth2.Token.RefreshToken
	token.Expiry = gaw.oAuth2.Token.Expiry

	authConf, _ := gads.NewCredentialsFromCode(context.TODO(), customerId, gaw.DeveloperToken, "Leapforce", &token)

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
		return nil, err
	}

	return &campaigns, nil
}
