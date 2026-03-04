package depsdev

import (
	"context"
	"dependency-dashboard/internal/domain"
	"dependency-dashboard/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	baseURL      = "https://api.deps.dev/v3"
	systemNPM    = "npm"
	relationRepo = "SOURCE_REPO"
)

type VersionResponse struct {
	VersionKey  VersionKey `json:"versionKey"`
	PublishedAt time.Time  `json:"publishedAt"`
	IsDefault   bool       `json:"isDefault"`
}

type VersionKey struct {
	System  string `json:"system"`
	Name    string `json:"name"`
	Version string `json:"version"`
}
type RelatedProjectsResponse struct {
	RelatedProjects []RelatedProject `json:"relatedProjects"`
}

type RelatedProject struct {
	ProjectKey         ProjectKey `json:"projectKey"`
	RelationProvenance string     `json:"relationProvenance"`
	RelationType       string     `json:"relationType"`
}

type ProjectKey struct {
	ID string `json:"id"`
}

type Client struct {
	baseURL string
	http    *retryablehttp.Client
}

func New(timeoutSec int) *Client {
	retryableHttpClient := retryablehttp.NewClient()
	retryableHttpClient.HTTPClient.Timeout = time.Duration(timeoutSec) * time.Second
	// TODO add cfg var for max retries
	retryableHttpClient.RetryMax = 3

	return &Client{
		baseURL: baseURL,
		http:    retryableHttpClient,
	}
}

func (c *Client) Fetch(ctx context.Context, name string) (*model.Dependency, error) {
	if name == "" {
		return nil, errors.New("package name is required")
	}

	version, err := c.fetchDefaultVersion(ctx, name)
	if err != nil {
		return nil, err
	}

	repo, err := c.fetchSourceRepo(ctx, name, version)
	if err != nil {
		return nil, err
	}

	score, err := c.fetchScore(ctx, repo)
	if err != nil {
		return nil, err
	}

	return &model.Dependency{
		Name:         name,
		Version:      version,
		OpenSSFScore: score,
		LastUpdated:  time.Now().UTC(),
	}, nil
}

func (c *Client) fetchDefaultVersion(ctx context.Context, name string) (string, error) {
	endpoint := fmt.Sprintf(
		"%s/systems/%s/packages/%s",
		c.baseURL,
		systemNPM,
		url.PathEscape(name),
	)

	var resp struct {
		Versions []VersionResponse `json:"versions"`
	}

	if err := c.get(ctx, endpoint, &resp); err != nil {
		return "", err
	}

	for _, v := range resp.Versions {
		if v.IsDefault {
			return v.VersionKey.Version, nil
		}
	}

	return "", errors.New("default version not found")
}

func (c *Client) fetchSourceRepo(ctx context.Context, name, version string) (string, error) {
	endpoint := fmt.Sprintf(
		"%s/systems/%s/packages/%s/versions/%s",
		c.baseURL,
		systemNPM,
		url.PathEscape(name),
		url.PathEscape(version),
	)

	var resp RelatedProjectsResponse
	if err := c.get(ctx, endpoint, &resp); err != nil {
		return "", err
	}

	for _, p := range resp.RelatedProjects {
		if p.RelationType == relationRepo {
			return p.ProjectKey.ID, nil
		}
	}

	return "", errors.New("source repository not found")
}

func (c *Client) fetchScore(ctx context.Context, projectID string) (float64, error) {
	if projectID == "" {
		return 0, errors.New("project id is empty")
	}

	endpoint := fmt.Sprintf(
		"%s/projects/%s",
		c.baseURL,
		url.PathEscape(projectID),
	)

	var resp struct {
		Scorecard struct {
			OverallScore float64 `json:"overallScore"`
		} `json:"scorecard"`
	}

	if err := c.get(ctx, endpoint, &resp); err != nil {
		return 0, err
	}

	return resp.Scorecard.OverallScore, nil
}

func (c *Client) get(ctx context.Context, url string, out any) error {
	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return domain.ErrNotFound
	}

	if resp.StatusCode >= 500 {
		return fmt.Errorf("upstream service error: %d", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		return domain.ErrUnexpected
	}

	return json.NewDecoder(resp.Body).Decode(out)
}
