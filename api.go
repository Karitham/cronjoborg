package cronjoborg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

const APIURL = "https://api.cron-job.org"

// Client is a cron-job.org client
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client
	apiKey string
}

func New(APIKey string, opts ...func(*Client)) *Client {
	c := &Client{
		apiKey: APIKey,
		client: &http.Client{Timeout: 5 * time.Second},
	}

	for _, f := range opts {
		f(c)
	}

	return c
}

func WithClient(httpc *http.Client) func(*Client) {
	return func(c *Client) {
		c.client = httpc
	}
}

// ListJobs
// https://docs.cron-job.org/rest-api.html#listing-cron-jobs
func (c *Client) ListJobs() ([]Job, error) {
	req, err := c.NewRequest(http.MethodGet, nil, "jobs")
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	type response struct {
		Jobs       []Job `json:"jobs"`
		SomeFailed bool  `json:"some_failed"`
	}

	var r response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return r.Jobs, nil
}

// GetJob
// https://docs.cron-job.org/rest-api.html#retrieving-cron-job-details
func (c *Client) GetJob(id int) (DetailedJob, error) {
	req, err := c.NewRequest(http.MethodGet, nil, "jobs", id)
	if err != nil {
		return DetailedJob{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return DetailedJob{}, err
	}
	defer resp.Body.Close()

	type response struct {
		JobDetails DetailedJob `json:"jobDetails"`
	}

	var r response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return DetailedJob{}, err
	}

	return r.JobDetails, nil
}

// CreateJob
// https://docs.cron-job.org/rest-api.html#creating-a-cron-job
func (c *Client) CreateJob(j DetailedJob) (int, error) {
	type request struct {
		Job DetailedJob `json:"job"`
	}
	reqPayload := request{Job: j}

	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(reqPayload); err != nil {
		return 0, err
	}

	req, err := c.NewRequest(http.MethodPut, body, "jobs")
	if err != nil {
		return 0, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	type response struct {
		JobID int `json:"jobId"`
	}

	var r response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return 0, err
	}

	return r.JobID, nil
}

// UpdateJob
// https://docs.cron-job.org/rest-api.html#updating-a-cron-job
func (c *Client) UpdateJob(id int, j DetailedJob) error {
	type request struct {
		Job DetailedJob `json:"job"`
	}
	reqPayload := request{Job: j}

	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(reqPayload); err != nil {
		return err
	}

	req, err := c.NewRequest(http.MethodPatch, body, "jobs", id)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// DeleteJob
// https://docs.cron-job.org/rest-api.html#deleting-a-cron-job
func (c *Client) DeleteJob(id int) error {
	req, err := c.NewRequest(http.MethodDelete, nil, "jobs", id)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// GetJobHistory
// https://docs.cron-job.org/rest-api.html#retrieving-the-job-execution-history
func (c *Client) GetJobHistory(id int) (JobHistory, error) {
	req, err := c.NewRequest(http.MethodGet, nil, "jobs", id, "history")
	if err != nil {
		return JobHistory{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return JobHistory{}, err
	}
	defer resp.Body.Close()

	var r JobHistory
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return JobHistory{}, err
	}

	return r, nil
}

// GetHistoryDetails
// https://docs.cron-job.org/rest-api.html#retrieving-job-execution-history-item-details
func (c *Client) GetHistoryDetails(jobID int, historyID int) (HistoryItem, error) {
	req, err := c.NewRequest(http.MethodGet, nil, "jobs", jobID, "history", historyID)
	if err != nil {
		return HistoryItem{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return HistoryItem{}, err
	}
	defer resp.Body.Close()

	type response struct {
		JobHistoryDetails HistoryItem `json:"jobHistoryDetails"`
	}
	var r response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return HistoryItem{}, err
	}

	return r.JobHistoryDetails, nil
}

func (c *Client) NewRequest(method string, body io.Reader, pathParts ...interface{}) (*http.Request, error) {
	parts := make([]string, 0, len(pathParts))
	for _, p := range pathParts {
		parts = append(parts, fmt.Sprint(p))
	}

	u, err := url.Parse(APIURL)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, path.Join(parts...))

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	c.headers(req)
	return req, nil
}

func (c *Client) headers(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")
}
