package plate

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

var (
	githubAPIBaseURL = "https://api.github.com"
	githubHTTPClient = &http.Client{Timeout: 15 * time.Second}
)

type githubRepositoryResponse struct {
	Private bool `json:"private"`
}

func (s *plateService) fetchKickplateYAML(ctx context.Context, repoURL, branch string, ownerAccountID string, organizationID *string) (*KickplateYAML, error) {
	return s.fetchKickplateYAMLWithOptions(ctx, repoURL, branch, ownerAccountID, organizationID, false)
}

func (s *plateService) fetchKickplateYAMLWithOptions(ctx context.Context, repoURL, branch string, ownerAccountID string, organizationID *string, forceRefresh bool) (*KickplateYAML, error) {
	manifest, err := s.fetchManifestYAML(ctx, repoURL, branch, "plate.yaml", ownerAccountID, organizationID, forceRefresh)
	if err == nil && strings.TrimSpace(manifest.Owner) != "" {
		return manifest, nil
	}

	if err != nil && !errors.Is(err, ErrMissingYAML) {
		return nil, err
	}

	legacy, legacyErr := s.fetchManifestYAML(ctx, repoURL, branch, "kikplate.yaml", ownerAccountID, organizationID, forceRefresh)
	if legacyErr != nil {
		if err == nil {
			return nil, ErrMissingYAML
		}
		if errors.Is(err, ErrMissingYAML) && errors.Is(legacyErr, ErrMissingYAML) {
			return nil, ErrMissingYAML
		}
		return nil, legacyErr
	}

	if strings.TrimSpace(legacy.Owner) == "" {
		return nil, ErrMissingYAML
	}

	return legacy, nil
}

func (s *plateService) fetchManifestYAML(ctx context.Context, repoURL, branch, filename string, ownerAccountID string, organizationID *string, forceRefresh bool) (*KickplateYAML, error) {
	apiURL := repoURLToContentsURL(repoURL, branch, filename)
	token, err := s.resolveGitHubToken(ctx, ownerAccountID, organizationID)
	if err != nil {
		return nil, err
	}
	resp, err := s.doGitHubRequest(apiURL, token, forceRefresh)
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

	var manifest KickplateYAML
	if err := yaml.Unmarshal(raw, &manifest); err != nil {
		return nil, ErrFetchFailed
	}

	return &manifest, nil
}

func (s *plateService) fetchRepositoryVisibility(ctx context.Context, repoURL string, ownerAccountID string, organizationID *string) (bool, error) {
	token, err := s.resolveGitHubToken(ctx, ownerAccountID, organizationID)
	if err != nil {
		return false, err
	}
	resp, err := s.doGitHubRequest(repoURLToAPIURL(repoURL), token, false)
	if err != nil {
		return false, ErrFetchFailed
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return false, fmt.Errorf("%w: github returned %d (%s)", ErrFetchFailed, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var repo githubRepositoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return false, ErrFetchFailed
	}

	return repo.Private, nil
}

func (s *plateService) doGitHubRequest(apiURL string, token string, forceRefresh bool) (*http.Response, error) {
	if forceRefresh {
		apiURL = fmt.Sprintf("%s&_nonce=%d", apiURL, time.Now().UnixNano())
	}

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "kickplate-api")
	if token = strings.TrimSpace(token); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}
	if forceRefresh {
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Pragma", "no-cache")
	}

	return githubHTTPClient.Do(req)
}

func (s *plateService) resolveGitHubToken(ctx context.Context, ownerAccountID string, organizationID *string) (string, error) {
	if s.githubApp == nil {
		return strings.TrimSpace(s.env.GitHubToken), nil
	}
	ownerUUID, err := uuid.Parse(ownerAccountID)
	if err != nil {
		return strings.TrimSpace(s.env.GitHubToken), nil
	}

	var orgUUID *uuid.UUID
	if organizationID != nil && strings.TrimSpace(*organizationID) != "" {
		parsed, parseErr := uuid.Parse(*organizationID)
		if parseErr == nil {
			orgUUID = &parsed
		}
	}

	token, err := s.githubApp.GetTokenForOwner(ctx, ownerUUID, orgUUID)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(token) != "" {
		return token, nil
	}
	return strings.TrimSpace(s.env.GitHubToken), nil
}

func privateRepositoryAllowed(repoPrivate bool, hasOrganization bool, privateOrganization bool) error {
	if repoPrivate && hasOrganization && !privateOrganization {
		return ErrPrivateRepositoryScope
	}
	return nil
}

func repoURLToContentsURL(repoURL, branch, filename string) string {
	return fmt.Sprintf("%s/repos/%s/contents/%s?ref=%s",
		strings.TrimRight(githubAPIBaseURL, "/"),
		extractRepoPath(repoURL), filename, url.QueryEscape(strings.TrimSpace(branch)))
}

func repoURLToAPIURL(repoURL string) string {
	return fmt.Sprintf("%s/repos/%s", strings.TrimRight(githubAPIBaseURL, "/"), extractRepoPath(repoURL))
}

func extractRepoPath(repoURL string) string {
	repoURL = strings.TrimSpace(repoURL)

	if strings.HasPrefix(repoURL, "git@github.com:") {
		repoURL = strings.TrimPrefix(repoURL, "git@github.com:")
	} else if strings.HasPrefix(repoURL, "ssh://git@github.com/") {
		repoURL = strings.TrimPrefix(repoURL, "ssh://git@github.com/")
	} else if strings.HasPrefix(repoURL, "https://github.com/") || strings.HasPrefix(repoURL, "http://github.com/") {
		if parsed, err := url.Parse(repoURL); err == nil {
			repoURL = strings.TrimPrefix(parsed.Path, "/")
		}
	}

	for _, prefix := range []string{
		"github.com/",
	} {
		if strings.HasPrefix(repoURL, prefix) {
			repoURL = strings.TrimPrefix(repoURL, prefix)
			break
		}
	}

	repoURL = strings.TrimSuffix(repoURL, ".git")
	repoURL = strings.SplitN(repoURL, "?", 2)[0]
	repoURL = strings.SplitN(repoURL, "#", 2)[0]
	repoURL = strings.Trim(repoURL, "/")

	segments := strings.Split(repoURL, "/")
	if len(segments) >= 2 {
		return path.Join(segments[0], segments[1])
	}

	return repoURL
}
