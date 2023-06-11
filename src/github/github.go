package github

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/yapaluc/hg-git/src/shell"

	"github.com/alessio/shellescape"
)

type PullRequest struct {
	BaseRefName string
	HeadRefName string
	State       string
	URL         string
	Title       string
	Body        string
}

func FetchPRForBranch(branchName string) (*PullRequest, error) {
	out, err := shell.Run(
		shell.Opt{},
		fmt.Sprintf(
			"gh pr list -s open -H %s --json url,state,title,baseRefName,headRefName,body",
			shellescape.Quote(branchName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("calling gh CLI: %w", err)
	}

	var resp []PullRequest
	err = json.Unmarshal([]byte(out), &resp)
	if err != nil {
		return nil, fmt.Errorf("decoding JSON from gh CLI: %w", err)
	}
	if len(resp) == 0 {
		// No PR.
		return nil, nil
	}
	pr := &resp[0]
	pr.Body = strings.ReplaceAll(pr.Body, "\r\n", "\n")
	return pr, nil
}

func FetchPRByURLOrNum(prURLOrNum string) (*PullRequest, error) {
	out, err := shell.Run(
		shell.Opt{},
		fmt.Sprintf(
			"gh pr view %s --json url,state,title,baseRefName,body",
			shellescape.Quote(prURLOrNum),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("calling gh CLI: %w", err)
	}

	var resp PullRequest
	err = json.Unmarshal([]byte(out), &resp)
	if err != nil {
		return nil, fmt.Errorf("decoding JSON from gh CLI: %w", err)
	}
	resp.Body = strings.ReplaceAll(resp.Body, "\r\n", "\n")
	return &resp, nil
}

func PRStrFromPRURL(prURL string) string {
	return fmt.Sprintf("#%d", PRNumFromPRURL(prURL))
}

func PRNumFromPRURL(prURL string) int {
	i, _ := strconv.Atoi(prURL[strings.LastIndex(prURL, "/")+1:])
	return i
}
