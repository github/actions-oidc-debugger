# actions-oidc-debugger

This action requests a JWT and prints the claims included within the JWT received from GitHub Actions.

## How to use this Action

Here's an example of how to use that action:

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
      - name: Debug OIDC Claims
        uses: github/actions-oidc-debugger@main
        with:
          audience: '${{ github.server_url }}/${{ github.repository_owner }}'
```
