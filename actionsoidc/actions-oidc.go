package actionsoidc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type ActionsOIDCClient struct {
	// the url to fetch the jwt
	TokenRequestURL string
	// the audience for the jwt
	Audience string
	// the token used to retrieve the jwt, not the jwt
	RequestToken string
}

type ActionsJWT struct {
	Count       int
	Value       string
	ParsedToken *jwt.Token
}

func GetEnvironmentVariable(e string) (string, error) {
	value := os.Getenv(e)
	if value == "" {
		return "", fmt.Errorf("missing %s from envrionment", e)
	}
	return value, nil
}

func QuitOnErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// construct a new ActionsOIDCClient
func NewActionsOIDCClient(tokenURL string, audience string, token string) (ActionsOIDCClient, error) {
	c := ActionsOIDCClient{
		TokenRequestURL: tokenURL,
		Audience:        audience,
		RequestToken:    token,
	}
	err := c.BuildTokenURL()
	return c, err
}

func DefaultOIDCClient(audience string) ActionsOIDCClient {
	tokenURL, err := GetEnvironmentVariable("ACTIONS_ID_TOKEN_REQUEST_URL")
	QuitOnErr(err)
	token, err := GetEnvironmentVariable("ACTIONS_ID_TOKEN_REQUEST_TOKEN")
	QuitOnErr(err)

	c, err := NewActionsOIDCClient(tokenURL, audience, token)
	QuitOnErr(err)

	return c
}

// this function uses an ActionsOIDCClient to build the complete URL
// to request a jwt
func (c *ActionsOIDCClient) BuildTokenURL() error {
	parsed_url, err := url.Parse(c.TokenRequestURL)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	if c.Audience != "" {
		query := parsed_url.Query()
		query.Set("audience", c.Audience)
		parsed_url.RawQuery = query.Encode()
		c.TokenRequestURL = parsed_url.String()
	}
	return nil
}

// retrieve an actions oidc token
func (c *ActionsOIDCClient) GetJWT() (*ActionsJWT, error) {
	request, err := http.NewRequest("GET", c.TokenRequestURL, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+c.RequestToken)

	var httpClient http.Client
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 from jwt api: %s", http.StatusText((response.StatusCode)))
	}

	rawBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var jwt ActionsJWT
	err = json.Unmarshal(rawBody, &jwt)

	return &jwt, err
}

func (j *ActionsJWT) Parse() {
	j.ParsedToken, _ = jwt.Parse(j.Value, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// we don't need a real check here
		return []byte{}, nil
	})
}

func (j *ActionsJWT) PrettyPrintClaims() string {
	if claims, ok := j.ParsedToken.Claims.(jwt.MapClaims); ok {
		jsonClaims, err := json.MarshalIndent(claims, "", "  ")
		if err != nil {
			fmt.Println(fmt.Errorf("%w", err))
		}
		return string(jsonClaims)
	}
	return ""
}
