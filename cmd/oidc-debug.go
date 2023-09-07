package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/github/actions-oidc-debugger/actionsoidc"
)

func main() {
	tokenPath := flag.String("token-path", "oidc-token", "the path to the token to debug")
	flag.Parse()

	if *tokenPath == "oidc-token" {
		actionsoidc.QuitOnErr(fmt.Errorf("token-path must be specified"))
	}

	c := actionsoidc.DefaultOIDCClient("blah-dont-care-right-now")

	contents, err := os.ReadFile(*tokenPath)
	actionsoidc.QuitOnErr(err)
	jwt, err := c.CreateOIDCClientFromValue(contents)
	actionsoidc.QuitOnErr(err)

	jwt.Parse()
	fmt.Print(jwt.PrettyPrintClaims())
}
