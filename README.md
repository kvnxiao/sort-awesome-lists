# sort-awesome-lists

Sort awesome lists by the number of stars in each GitHub repository, for each sub-heading / section in the list.

For example:
  - [**awesome-go** sorted by number of stars](https://gist.github.com/kvnxiao/cb432fca8cd9b59e325286b8f33cf53d)
  - [**awesome-rust** sorted by number of stars](https://gist.github.com/kvnxiao/fe8cd6ca03978a2ee69e36a37251bcd2)
  - [**awesome-kotlin** sorted by number of stars](https://gist.github.com/kvnxiao/5f809440525304c918b553b4bbc8cd73)
  - [**awesome-java** sorted by number of stars](https://gist.github.com/kvnxiao/dfea78544dd74953453ba74f6e59ee6f)

This is a CLI application written in Go and uses 0 external dependencies.

Sorting by stars implies parsing the original awesome-list `README.md` file and outputting a modified version where each section is sorted in descending order by the number of stars (for each valid github repository).

GitHub repository detection involves checking each markdown bullet point if it contains a `username.github.io/repo` or `github.com/username/repo` link. Otherwise, if a project website is linked, the application will attempt to download and parse the webpage to check if a GitHub repository link exists within the HTML.

## How to use

`sort-awesome-lists` is a CLI application. Build it and run in your terminal.

### Building

```
go build -o sort-awesome-lists main.go
```

Creates an executable file called `sort-awesome-lists` in your directory. Run in your terminal with `./sort-awesome-lists`

### Usage

```
Usage of sort-awesome-lists:
  -bs int
        number of concurrent requests to send to GitHub API at a time, per each block found. (default 5)
  -o string
        name of file to write output to if set, otherwise prints to stdout
  -t string
        GitHub personal access token
  -v    prints debug messages to stdout if true (default = false)
```

A GitHub personal access token is **required** by the `-t` flag as this CLI application hits the GitHub API for repository statistics. The token allows one to access the GitHub API at a rate-limit of 5000 requests per hour. A personal access token with 0 permissions checked can be generated and used (go here to create one if you don't already have one: https://github.com/settings/tokens)

This tool currently supports `username.github.io/repo` and `github.com/username/repo` detection.

#### Example:

```
./sort-awesome-lists -t="$token" -o="awesome-go-sorted.md" https://raw.githubusercontent.com/avelino/awesome-go/master/README.md
```
where `$token` is your github personal access token.

The above example will download and parse the markdown file from `https://raw.githubusercontent.com/avelino/awesome-go/master/README.md`, and output a sorted markdown output in a file called `awesome-go-sorted.md` in the same working directory.

### Known Issues / Gotchas

For entries in a list that do not directly link to `github.com/username/repo` or `username.github.io/repo`, the webpage will be downloaded and parsed to check if a GitHub repository link exists within the HTML. This means that this tool will be unable to pick up any links for websites that use a JavaScript framework and require JavaScript to render the page.
