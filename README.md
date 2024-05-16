# actions-oidc-debugger

This action requests a JWT and prints the claims included within the JWT received from GitHub Actions.

## How to use this Action

To use this Action in another repository, you must checkout this Action repo and then run it.
Here's an example of how that is done:

```yaml

on: [pull_request]

jobs:
  oidc_debug_test:
    permissions:
      contents: read
      id-token: write
    runs-on: ubuntu-latest
    name: A test of the oidc debugger
    steps:
      - name: Checkout actions-oidc-debugger
        uses: actions/checkout@v3
        with:
          repository: github/actions-oidc-debugger
          ref: main
          token: ${{ secrets.your-checkout-token }}
          path: ./.github/actions/actions-oidc-debugger
      - name: Debug OIDC Claims
        uses: ./.github/actions/actions-oidc-debugger
        with:
          audience: '${{ github.server_url }}/${{ github.repository_owner }}'
```
