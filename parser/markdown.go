package parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kvnxiao/sort-awesome-lists/github"
	"github.com/kvnxiao/sort-awesome-lists/logging"
	"github.com/kvnxiao/sort-awesome-lists/requests"
)

var (
	rLine = regexp.MustCompile(`^\s*([*\-]) \[.*?]\((https*|mailto):`)
	rUrl  = regexp.MustCompile(`\((https*://.*?)\)`)
)

type Repository struct {
	url       *url.URL
	text      string
	stars     int
	repoURL   string
	separator string
}

type GithubBlock struct {
	start        int
	end          int
	repositories []*Repository
}

type ByStars []*Repository

func (s ByStars) Len() int {
	return len(s)
}

func (s ByStars) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByStars) Less(i, j int) bool {
	ri := s[i]
	rj := s[j]
	if ri.stars == rj.stars {
		// sort ascending on lexicographical string order
		return ri.repoURL < rj.repoURL
	} else {
		// sort descending on stars
		return ri.stars > rj.stars
	}
}

type Markdown struct {
	lines  []string
	blocks []*GithubBlock
}

func ParseMarkdown(url string) *Markdown {
	logging.Println("Retrieving markdown...")
	now := time.Now()
	resp, err := requests.Get(url, nil)
	if err != nil {
		log.Fatalf("an error occurred retrieving markdown: %v", err)
	}
	defer resp.Body.Close()
	took := time.Now().Sub(now)
	logging.Printlnf("Markdown retrieved in %v", took.String())

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("couldn't read response body: %v", err)
	}

	markdownBody := string(b)
	lines := strings.Split(markdownBody, "\n")

	marked := false
	var blocks []*GithubBlock
	var repositories []*Repository
	start := 0
	end := 0
	for i, line := range lines {
		submatches := rLine.FindStringSubmatch(line)
		if len(submatches) > 0 {
			separator := submatches[1]
			if !marked {
				marked = true
				start = i
				end = i
			} else {
				end++
			}
			repositories = append(repositories, parseRepoText(line, separator))
		} else {
			if marked {
				blocks = append(blocks, &GithubBlock{
					start:        start,
					end:          end,
					repositories: repositories,
				})
				repositories = nil
			}
			marked = false
		}
	}
	if marked {
		blocks = append(blocks, &GithubBlock{
			start:        start,
			end:          end,
			repositories: repositories,
		})
		repositories = nil
	}
	return &Markdown{
		lines:  lines,
		blocks: blocks,
	}
}

func parseRepoText(line, separator string) *Repository {
	submatch := rUrl.FindStringSubmatch(line)
	if len(submatch) < 2 {
		return &Repository{
			text:    line,
			url:     nil,
			stars:   0,
			repoURL: "",
			separator: separator,
		}
	}

	urlString := submatch[1]
	u, err := url.Parse(urlString)
	if err != nil {
		log.Fatalf("an error occurred parsing repositories %s: %v", urlString, err)
	}

	// parse hostname and path for potential github repo api endpoint
	hostname := u.Hostname()
	path := u.Path
	repoURL := github.GetApiEndpoint(hostname, path)

	return &Repository{
		text:    line,
		url:     u,
		stars:   0,
		repoURL: repoURL,
		separator: separator,
	}
}

func (md *Markdown) CountAll() int {
	c := 0
	for _, block := range md.blocks {
		c += len(block.repositories)
	}
	return c
}

func (md *Markdown) FetchStars(token string) {
	blockCount := len(md.blocks)

	logging.Printlnf("%d blocks to fetch info for", blockCount)
	for i, githubBlock := range md.blocks {
		githubBlock.fetchStars(token, i)
	}
}

func (md *Markdown) Sort() {
	for blockNum, githubBlock := range md.blocks {
		logging.Printlnf("Sorting block %d", blockNum)
		sort.Sort(ByStars(githubBlock.repositories))

		start := githubBlock.start
		for i, repo := range githubBlock.repositories {
			index := start + i
			numStr := strings.Replace(fmt.Sprintf("<code>%6s</code>", strconv.Itoa(repo.stars)), " ", "&nbsp;", -1)
			indexOfFirstSeparator := strings.Index(repo.text, repo.separator + " ")
			md.lines[index] = repo.text[:indexOfFirstSeparator] + repo.separator + " **" + numStr + "** " + repo.text[indexOfFirstSeparator+2:]
		}
	}
}

func (md *Markdown) ToString() string {
	return strings.Join(md.lines, "\n")
}

func (b *GithubBlock) fetchStars(token string, blockNumber int) {
	repoCount := len(b.repositories)
	var wg sync.WaitGroup
	wg.Add(repoCount)

	logging.Printlnf("Started fetching stars for block %d.", blockNumber)
	for _, repository := range b.repositories {
		repository := repository

		go func(repo *Repository) {
			if repo.repoURL != "" {
				repo.stars = github.GetRepoStars(repository.repoURL, token)
			} else {
				repo.stars = 0
			}

			wg.Done()
		}(repository)

	}
	wg.Wait()
	logging.Printlnf("fetching stars for block %d done.", blockNumber)
}
