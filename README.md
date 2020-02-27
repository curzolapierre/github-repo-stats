# github-repo-stats

## How to Build?

The project is using Go, then you need a running Go environment: [Official documentation](https://golang.org/doc/install)

Once that's done, all you have to do is to `go get` the project, with the following command:

```shell
go get github.com/curzolapierre/github-repo-stats
```

That's it, you've build the latest version of the github repository project (the binary will be present in `$GOPATH/bin/github-repo-stats`)

To fully exploit the project, you'll need to create the config file `./server.config.json` and set the `personal_token` value with a `Personal access tokens` (you can create it [here](https://github.com/settings/tokens)). Use `./config.server.example.json` as template.

## Rate limiting

- 60 requests per hour for unauthenticated requests
- 5000 requests per hour for authenticated requests

if rate limit is exceed, github API will return 403 Forbidden

## OAuth2 token

Github documentation to help you to create a personal tokens: [here](https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line)

Need to be set in request header: `Authorization: token OAUTH-TOKEN`

## endpoints

Root endpoint: `https://api.github.com`

### list public repositories

Lists all public repositories in the order that they were created.

GET `/repositories`

Response:

Status: 200 OK

```JSON
[
  {
    "id": 1296269,
    "node_id": "MDEwOlJlcG9zaXRvcnkxMjk2MjY5",
    "name": "Hello-World",
    "full_name": "octocat/Hello-World",
    "...": "..."
  }
]
```

### list languages

Lists languages for the specified repository. The value shown for each language is the number of bytes of code written in that language

GET `/repos/:owner/:repo/languages`

Response:

Status: 200 OK

```JSON
{
  "HTML": 76092,
  "Ruby": 24144,
  "JavaScript": 14319,
  "CSS": 12692,
  "Dockerfile": 1284
}
```
