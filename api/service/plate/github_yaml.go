package plate

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"gopkg.in/yaml.v3"
)

func (s *plateService) fetchKickplateYAML(repoURL, branch string) (*KickplateYAML, error) {
	apiURL := repoURLToContentsURL(repoURL, branch)
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, ErrFetchFailed
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "kickplate-api")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, ErrFetchFailed
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrMissingYAML
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("%w: github returned %d (%s)", ErrFetchFailed, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var ghResp struct {
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ghResp); err != nil {
		return nil, ErrFetchFailed
	}

	var raw []byte
	if ghResp.Encoding == "base64" {
		raw, err = base64.StdEncoding.DecodeString(ghResp.Content)
		if err != nil {
			return nil, ErrFetchFailed
		}
	} else {
		raw = []byte(ghResp.Content)
	}

	var kp KickplateYAML
	if err := yaml.Unmarshal(raw, &kp); err != nil {
		return nil, ErrFetchFailed
	}
	if kp.Owner == "" {
		return nil, ErrMissingYAML
	}

	return &kp, nil
}

func repoURLToContentsURL(repoURL, branch string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/contents/kikplate.yaml?ref=%s",
		extractRepoPath(repoURL), url.QueryEscape(strings.TrimSpace(branch)))
}

func extractRepoPath(repoURL string) string {
	repoURL = strings.TrimSpace(repoURL)

	repoURL = strings.TrimPrefix(repoURL, "git@github.com:")

	for _, prefix := range []string{
		"https://github.com/",
		"http://github.com/",
		"github.com/",
	} {
		if strings.HasPrefix(repoURL, prefix) {
			repoURL = strings.TrimPrefix(repoURL, prefix)
			break
		}
	}

	repoURL = strings.TrimSuffix(repoURL, ".git")
	repoURL = strings.Trim(repoURL, "/")
	return repoURL
}
