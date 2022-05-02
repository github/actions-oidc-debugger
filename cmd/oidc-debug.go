package main

import (
	"flag"
	"fmt"

	"github.com/github/actions-oidc-debugger/actionsoidc"
)

func main() {

	audience := flag.String("audience", "https://github.com/", "the audience for the requested jwt")
	flag.Parse()

	if *audience == "https://github.com/" {
		actionsoidc.QuitOnErr(fmt.Errorf("-audience cli argument must be specified"))
	}

	c := actionsoidc.DefaultOIDCClient(*audience)
	jwt, err := c.GetJWT()
	actionsoidc.QuitOnErr(err)

	jwt.Parse()
	fmt.Print(jwt.PrettyPrintClaims())
}
