package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type help_t struct {
	name string
	text string
}

var CommandHelp = []help_t{
	{"help", "display this help text"},
	{"list", "list all repositories"},
	{"find", "search go package git repositories for keywords"},
}

func Github() *github.Client {
	client := github.NewClient(nil)
	return client
}

func PrintHelp() {
	for _, helptext := range CommandHelp {
		fmt.Fprintf(os.Stderr, "%s - %s\n", helptext.name, helptext.text)
	}
	fmt.Fprintf(os.Stderr, "\n")
}

func main() {
	var accessToken = flag.String("token", os.Getenv("GITHUB_TOKEN"), "token github api")
	var user = flag.String("user", "", "user to list repos for")
	var allLanguages = flag.Bool("all", false, "dont return only go language repos")
	flag.Usage = func() {
		PrintHelp()
		flag.PrintDefaults()
	}
	flag.Parse()

	var client = getclient(*accessToken)
	switch flag.Arg(0) {
	case "", "help":
		flag.Usage()
	case "list":
		listRepositories(client, *user)
	case "code":
		if flag.NArg() == 1 {
			flag.Usage()
			return
		}
		findCodeRepositories(client, strings.Join(flag.Args()[1:], "+"), *allLanguages)
	case "find":
		if flag.NArg() == 1 {
			flag.Usage()
			return
		}
		findRepositories(client, strings.Join(flag.Args()[1:], "+"), *allLanguages)
	}

}

func getclient(token string) *github.Client {
	var client *github.Client
	if token != "" {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)

		client = github.NewClient(tc)

	} else {
		client = Github()
	}

	return client
}

func listRepositories(client *github.Client, name string) {
	// list all repositories for the authenticated user
	ctx := context.Background()
	for i := range []int{1, 2} { // pages of 100
		opts := &github.RepositoryListOptions{}
		opts.PerPage = 100
		opts.Page = i
		repos, _, err := client.Repositories.List(ctx, name, opts)
		if err != nil {
			fmt.Println(err)
			os.Exit(111)
		}
		if len(repos) == 0 {
			break
		}
		for _, repo := range repos {
			descr := repo.GetDescription()
			if descr != "" {
				descr = " - " + descr
			}
			fmt.Printf("%s/%s%s\n", name, repo.GetName(), descr)

		}
	}
}
func findCodeRepositories(client *github.Client, keywords string, allLanguages bool) {
	// log.Printf("searching repos for: %q", keywords)
	if !allLanguages {
		keywords += "+language:go"
	}
	ctx := context.Background()
	x := 0
	for i := range []int{1, 2} { // pages of 100
		opts := &github.SearchOptions{}
		opts.PerPage = 100
		opts.Page = i

		repos, _, err := client.Search.Code(ctx, keywords, opts)
		if err != nil {
			fmt.Println(err)
			os.Exit(111)
		}
		num := repos.GetTotal()
		fmt.Printf("Found %v results\n", num)
		fmt.Printf("# (*) repository - description\n")
		//already := map[string]bool{}
		for _, repocode := range repos.CodeResults {
			repo := repocode.GetRepository()
			fmt.Println(repocode.GetHTMLURL())
			descr := repo.GetDescription()
			if descr != "" {
				descr = " - " + descr
			}
			stars := repo.GetStargazersCount()
			fmt.Printf("%00v (%0v) %s%s\n", x, stars, repo.GetCloneURL(), descr)
			x++
		}
		if num != 100 {
			break
		}
	}
}
func findRepositories(client *github.Client, keywords string, allLanguages bool) {
	// log.Printf("searching repos for: %q", keywords)
	if !allLanguages {
		keywords += "+language:go"
	}
	ctx := context.Background()
	x := 0
	for i := range []int{1, 2} { // pages of 100
		opts := &github.SearchOptions{}
		opts.PerPage = 100
		opts.Page = i

		repos, _, err := client.Search.Repositories(ctx, keywords, opts)
		if err != nil {
			fmt.Println(err)
			os.Exit(111)
		}
		num := repos.GetTotal()
		fmt.Printf("Found %v results\n", num)
		fmt.Printf("# (*) repository - description\n")
		for _, repo := range repos.Repositories {
			descr := repo.GetDescription()
			if descr != "" {
				descr = " - " + descr
			}
			stars := repo.GetStargazersCount()
			fmt.Printf("%00v (%0v) %s%s\n", x, stars, repo.GetCloneURL(), descr)
			x++

		}
		if num != 100 {
			break
		}
	}
}
