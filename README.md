# actions-oidc-debugger

This action requests a JWT and prints the claims included within the JWT received from GitHub Actions.

## Usage

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
      - name: Checkout
        uses: actions/checkout@v3
      - name: Debug OIDC Claims
        uses: github/actions-oidc-debugger@v1
        with:
          audience: 'https://github.com/github
```
