package parser

import (
	"regexp"
	"sort"
	"strings"
)

// Result is a Parse result, returning the matched repo, issue, etc. as applicable
type Result struct {
	User  string // the matched user, if applicable
	Match string // the matched shorthand value, if applicable
	Issue string // the matched issue number, if applicable
	Path  string // the matched path fragment, if applicable
	Query string // the remainder of the input, if not otherwise parsed
	owner string
	name  string
}

func (r *Result) HasRepo() bool {
	return len(r.name) > 0
}

func (r *Result) SetRepo(repo string) error {
	parts := strings.SplitN(repo, "/", 2)
	if len(parts) > 1 {
		r.owner = parts[0]
		r.name = parts[1]
	}
	return nil
}

func (r *Result) Repo() string {
	if r.HasRepo() {
		return r.owner + "/" + r.name
	}
	return ""
}

// Annotation is a helper for displaying details about a match. Returns a string
// with a leading space, noting the matched shorthand and issue if applicable.
func (r *Result) Annotation() (ann string) {
	if len(r.Match) > 0 {
		ann += " (" + r.Match
		if len(r.Issue) > 0 {
			ann += "#" + r.Issue
		}
		ann += ")"
	}
	return
}

// Parse takes a repo mapping and input string and attempts to extract a repo,
// issue, etc. from the input using the repo map for shorthand expansion.
func Parse(repoMap, userMap map[string]string, input string) *Result {
	path := ""
	user := ""
	owner, name, match, query := extractRepo(repoMap, input)
	if len(name) == 0 {
		owner, match, query = extractUser(userMap, input)
		user = owner
	}
	issue, query := extractIssue(query)
	if issue == "" {
		path, query = extractPath(query)
	}
	return &Result{
		owner: owner,
		name:  name,
		User:  user,
		Match: match,
		Issue: issue,
		Path:  path,
		Query: query,
	}
}

var (
	userRepoRegexp = regexp.MustCompile(`^([A-Za-z0-9][-A-Za-z0-9]*)/([\w\.\-]+)\b`) // user/repo
	issueRegexp    = regexp.MustCompile(`^#?([1-9]\d*)$`)
	pathRegexp     = regexp.MustCompile(`^(/\S*)$`)
)

func extractRepo(repoMap map[string]string, input string) (owner, name, match, query string) {
	var keys []string
	for k := range repoMap {
		keys = append(keys, k)
	}

	// sort the keys in reverse so the longest is matched first
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	for _, k := range keys {
		if strings.HasPrefix(input, k) {
			parts := strings.SplitN(repoMap[k], "/", 2)
			if len(parts) > 1 {
				return parts[0], parts[1], k, strings.TrimLeft(input[len(k):], " ")
			}
		}
	}

	result := userRepoRegexp.FindStringSubmatch(input)
	if len(result) > 0 {
		repo, owner, name := result[0], result[1], result[2]
		return owner, name, "", strings.TrimLeft(input[len(repo):], " ")
	}
	return "", "", "", input
}

func extractUser(userMap map[string]string, input string) (user, match, query string) {
	var keys []string
	for k := range userMap {
		keys = append(keys, k)
	}

	// sort the keys in reverse so the longest is matched first
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	for _, k := range keys {
		if strings.HasPrefix(input, k) {
			return userMap[k], k, strings.TrimLeft(input[len(k):], " ")
		}
	}

	return "", "", input
}

func extractIssue(query string) (issue, remainder string) {
	match := issueRegexp.FindStringSubmatch(query)
	if len(match) > 0 {
		issue = match[1]
		remainder = ""
	} else {
		issue = ""
		remainder = query
	}
	return
}

func extractPath(query string) (path, remainder string) {
	match := pathRegexp.FindStringSubmatch(query)
	if len(match) > 0 {
		path = match[1]
		remainder = ""
	} else {
		path = ""
		remainder = query
	}
	return
}
