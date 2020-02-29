# github-repo-stats

## How to Build?

The project is using Go, then you need a running Go environment: [Official documentation](https://golang.org/doc/install)

Once that's done, all you have to do is to `go get` the project, with the following command:

```shell
go get github.com/curzolapierre/github-repo-stats
```

That's it, you've build the latest version of the github repository project (the binary will be present in `$GOPATH/bin/github-repo-stats`)

To fully exploit the project, you'll need to create the config file `./credentials.json` and set the `personal_token` value with a `Personal access tokens` (you can create it [here](https://github.com/settings/tokens)). Use `./credentials.example.json` as template.

## How it works

First a request will be made to fetch repositories according to user input

If user input is empty last repositories from current day - 1 will be fetched

For each repositories, a new request will be made to retrieve its languages

This process will be made by workers, where their number can be set in `./server.config.json` with the `worker_number` field

## Rate limiting

- 60 requests per hour for unauthenticated requests
- 5000 requests per hour for authenticated requests

## OAuth2 token

Github documentation to help you to create a personal tokens: [here](https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line)

Need to be set in request header: `Authorization: token OAUTH-TOKEN`

## endpoints

Root endpoint: `https://api.github.com`

### list public repositories

[Official documentation](https://developer.github.com/v3/search/#search-repositories)
Lists all public repositories in the order that they were created.

GET `/search/repositories`

Response:

Status: 200 OK

```JSON
{
  "total_count": 123,
  "incomplete_results": true,
  "items": [
    {
      "id": 1296269,
      "node_id": "MDEwOlJlcG9zaXRvcnkxMjk2MjY5",
      "name": "Hello-World",
      "full_name": "octocat/Hello-World",
      "...": "..."
    },
  ]
```

### list languages

[Official documentation](https://developer.github.com/v3/repos/#list-languages)
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
