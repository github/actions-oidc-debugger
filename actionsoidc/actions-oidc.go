package actionsoidc

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"

	"github.com/golang-jwt/jwt"
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

type JWK struct {
	N   string
	Kty string
	Kid string
	Alg string
	E   string
	Use string
	X5c []string
	X5t string
}

type JWKS struct {
	Keys []JWK
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

func (c *ActionsOIDCClient) CreateOIDCClientFromValue(rawValue []byte) (*ActionsJWT, error) {
	var jwt ActionsJWT
	jwt.Value = string(rawValue)
	return &jwt, nil
}

func getKeyFromJwks(jwksBytes []byte) func(*jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		var jwks JWKS
		if err := json.Unmarshal(jwksBytes, &jwks); err != nil {
			return nil, fmt.Errorf("Unable to parse JWKS")
		}

		for _, jwk := range jwks.Keys {
			if jwk.Kid == token.Header["kid"] {
				nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
				if err != nil {
					return nil, fmt.Errorf("Unable to parse key")
				}
				var n big.Int

				eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
				if err != nil {
					return nil, fmt.Errorf("Unable to parse key")
				}
				var e big.Int

				key := rsa.PublicKey{
					N: n.SetBytes(nBytes),
					E: int(e.SetBytes(eBytes).Uint64()),
				}

				return &key, nil
			}
		}

		return nil, fmt.Errorf("Unknown kid: %v", token.Header["kid"])
	}
}

func (j *ActionsJWT) Parse() {
	// get jwks
	resp, err := http.Get("https://token.actions.githubusercontent.com/.well-known/jwks")
	if err != nil {
		fmt.Println(err)
	}

	jwksBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	token, err := jwt.Parse(string(j.Value), getKeyFromJwks(jwksBytes))
	if err != nil || !token.Valid {
		fmt.Println("unable to validate jwt")
		log.Fatal(err)
	}

	// token, err := jwt.Parse(j.Value, func(token *jwt.Token) (interface{}, error) {
	// 	// Don't forget to validate the alg is what you expect:
	// 	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
	// 		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	// 	}

	// we don't need a real check here
	// return []byte{}, nil
	// })

	j.ParsedToken = token
}

func (j *ActionsJWT) PrettyPrintClaims() string {
	fmt.Println(j.ParsedToken.Header)

	if claims, ok := j.ParsedToken.Claims.(jwt.MapClaims); ok {
		jsonClaims, err := json.MarshalIndent(claims, "", "  ")
		if err != nil {
			fmt.Println(fmt.Errorf("%w", err))
		}
		return string(jsonClaims)
	}
	return ""
}
