# actions-oidc-debugger

This action requests a JWT and prints the claims included within the JWT received from GitHub Actions.

## How to use this Action

Here's an example of how to use this action:

```yaml

name: Test Debugger Action
on: 
  pull_request:
  workflow_dispatch:

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

The resulting output in your Actions log will look something like this:

```json
{
  "actor": "GrantBirki",
  "actor_id": "23362539",
  "aud": "https://github.com/github",
  "base_ref": "main",
  "enterprise": "github",
  "enterprise_id": "11468",
  "event_name": "pull_request",
  "exp": 1751581975,
  "head_ref": "release-setup",
  "iat": 1751560375,
  "iss": "https://token.actions.githubusercontent.com",
  "job_workflow_ref": "github/actions-oidc-debugger/.github/workflows/action-test.yml@refs/pull/27/merge",
  "job_workflow_sha": "7f93a73b8273af5d35fcd70661704c1cadc57054",
  "jti": "4a576b35-ff09-41c5-af2c-ca62dd89b76a",
  "nbf": 1751560075,
  "ref": "refs/pull/27/merge",
  "ref_protected": "false",
  "ref_type": "branch",
  "repository": "github/actions-oidc-debugger",
  "repository_id": "487920697",
  "repository_owner": "github",
  "repository_owner_id": "9919",
  "repository_visibility": "public",
  "run_attempt": "1",
  "run_id": "16055869479",
  "run_number": "33",
  "runner_environment": "github-hosted",
  "sha": "7f93a73b8273af5d35fcd70661704c1cadc57054",
  "sub": "repo:github/actions-oidc-debugger:pull_request",
  "workflow": "Test Debugger Action",
  "workflow_ref": "github/actions-oidc-debugger/.github/workflows/action-test.yml@refs/pull/27/merge",
  "workflow_sha": "7f93a73b8273af5d35fcd70661704c1cadc57054"
}
```
