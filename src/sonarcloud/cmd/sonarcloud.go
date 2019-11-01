// Package sonarcloud API wrapper
//
// Provides a basic wrapper for the SonarCloud API
//   Support is limited to:
//		projects
//		tokens
//		metrics
//		qaulity gates
package sonarcloud

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ProjectSearchResponse for a GET / search on projects
type ProjectSearchResponse struct {
	// Paging object
	Paging PagingObject `json:"paging"`

	// Components list
	Components []ComponentsObject `json:"components"`
}

// PagingObject used for response that are limited by size
type PagingObject struct {

	// Index page number
	Index int `json:"pageIndex"`

	// Size elements on this page
	Size int `json:"pageSize"`

	// Total number of pages
	Total int `json:"total"`
}

// ComponentsObject structure
type ComponentsObject struct {
	// Organization name
	Organization string `json:"organization"`

	// Key for access this object
	Key string `json:"key"`

	// Name for display
	Name string `json:"name"`

	// Qualifier is type of component
	Qualifier string `json:"qualifier"`

	// LastAnalysisDate time of last CI
	LastAnalysisDate string `json:"lastAnalysisDate"`

	// Revision hash
	Revision string `json:"revision"`
}

// NewProject Used to create a new project
//
type NewProject struct {
	// Organization (required) is a valid SonarCloud organization
	Organization string `json:"organization"`

	// Name friendly name for display
	Name string `json:"name"`

	// Project is the SonarCloud Key
	Project string `json:"project"`

	// Visibility (optional) private or public
	Visibility string `json:"visibility"`
}

// NewProjectResponse includes wrapper "project"
// This sucks we can't use the same structure but they
// change field names
type NewProjectResponse struct {
	Project NewProjectResponseObject `json:"project"`
}

// NewProjectResponseObject Object include in response
//
type NewProjectResponseObject struct {
	// Key
	Key string `json:"key"`

	// Name friendly name for display
	Name string `json:"name"`

	// Qualifier is the SonarCloud component type
	Qualifier string `json:"qualifier"`

	// Visibility (optional) private or public
	Visibility string `json:"visibility"`
}

// sonarcloudclient
//   type and methods used for accessing SonarCloud API
type sonarcloudclient struct {
	//   Client is an http.Client created when New() is called
	Client *http.Client

	//   Host is the default host, sonarcloud.io by default
	Host string

	//   APIversion is the api prefix to use in API calls /api
	APIVersion string

	//   Token used to authenticate
	Token string

	// connectino string
	URI string
}

// NewTokenResponse holds response from user_tokens/generate
type NewTokenResponse struct {
	// Login name token is for
	Login string `json:"login"`

	// Name of the token
	Name string `json:"name"`

	// The Token
	Token string `json:"token"`

	// CreatedAt date and time of creation
	CreatedAt string `json:"createdAt"`
}

// GetTokenResponse user_tokens/search returns a user and
// a list of their tokens
type GetTokenResponse struct {
	// Login name of user
	Login string `json:"login"`

	// List of tokens
	Tokens []GetTokenItem `json:"userTokens"`
}

// GetTokenItem items returned in a token search
type GetTokenItem struct {
	// Name of the token
	Name string `json:"name"`

	// CreatedAt date and time
	CreatedAt string `json:"createdAt"`

	// LastConnectionDate date and time token was last used
	// Only updated hourly
	LastConnectionDate string `json:"lastConnectionDate"`
}

// SonarCloudError
type sonarCloudError struct {
	errNumber int
	errMsg    string
}

func (e *sonarCloudError) Error() string {
	msg := fmt.Sprintf("Err: %v, %v\n", e.errNumber, e.errMsg)
	return msg
}

// New(sondarcloudclient, token)
//   token is a valid sonarcloud user token
//	 if must have admin access
func (c *sonarcloudclient) New(token string) error {

	c.Client = &http.Client{
		Timeout: 10 * time.Second,
	}

	if c.Host == "" {
		c.Host = DefaultHost
	}

	if c.APIVersion == "" {
		c.APIVersion = DefaultAPI
	}

	if token == "" {
		a := sonarCloudError{errNumber: -1, errMsg: "Token is require"}
		fmt.Println(a.Error())
		return nil
	}

	c.Token = token

	// https://token@host
	c.URI = fmt.Sprintf("%s%s@%s", DefaultScheme, c.Token, c.Host)

	return nil
}

// GetProject
// TODO make this a vardic fucntion taking a list of project names
func (c *sonarcloudclient) GetProject(org, name string) (*http.Response, error) {
	options := "?"
	options += fmt.Sprintf(Projects, name)
	options += fmt.Sprintf("&"+Organization, org)

	url := c.URI + ProjectSearch + options
	req, err := http.NewRequest("GET", url, nil)
	resp, err := c.Client.Do(req)

	if err != nil {
		fmt.Println("Error is ", err)
		return resp, err
	}

	return resp, nil
}

// Create Project(Project)
//   Create a new SonarCloud project usinig p
//   Note SonarCloud expects application/x-www-form-urlencoded
func (c *sonarcloudclient) CreateProject(p NewProject) (*http.Response, error) {

	data := url.Values{}
	data.Set("name", p.Name)

	// Project names for none default organization are global
	// Add a prefix to avoid naming conflicts
	// TODO: Make this configurable by the end user
	data.Set("project", KeyPrefix+p.Project)
	data.Set("organization", p.Organization)
	data.Set("visibility", p.Visibility)

	url := c.URI + ProjectCreate
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	rsp, err := c.Client.Do(req)

	if rsp.StatusCode == 400 {
		errmsg, _ := ioutil.ReadAll(rsp.Body)
		return rsp, errors.New(string(errmsg))
	}

	if err != nil {
		fmt.Println("Error is ", err)
		return rsp, err
	}
	return rsp, nil
}

// Delete Project(projectKey)
//   Delete a new SonarCloud project usinig p
//   Note SonarCloud expects application/x-www-form-urlencoded
func (c *sonarcloudclient) DeleteProject(p string) (*http.Response, error) {

	// Project names for none default organization are global
	pk := KeyPrefix + p
	data := url.Values{}
	data.Set("project", pk)

	url := c.URI + ProjectDelete
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	rsp, err := c.Client.Do(req)

	if rsp.StatusCode == 400 {
		errmsg, _ := ioutil.ReadAll(rsp.Body)
		return rsp, errors.New(string(errmsg))
	}

	if err != nil {
		fmt.Println("Error is ", err)
		return rsp, err
	}
	return rsp, nil
}

// Create Token(tn string)
//   Create a new SonarCloud token with the name tn
//   Note SonarCloud expects application/x-www-form-urlencoded
func (c *sonarcloudclient) CreateToken(tn string) (*http.Response, error) {

	data := url.Values{}
	data.Set("name", tn)

	url := c.URI + TokenCreate
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	rsp, err := c.Client.Do(req)

	// return errmsg in the body as the error
	if rsp.StatusCode == 400 {
		errmsg, _ := ioutil.ReadAll(rsp.Body)
		return rsp, errors.New(string(errmsg))
	}

	if err != nil {
		fmt.Println("Error is ", err)
		return rsp, err
	}
	return rsp, nil
}

// Revoke Token(tn string)
//   Revoke a SonarCloud token with the name tn
func (c *sonarcloudclient) RevokeToken(tn string) (*http.Response, error) {

	data := url.Values{}
	data.Set("name", tn)

	url := c.URI + TokenRevoke
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	rsp, err := c.Client.Do(req)

	// There was a problem with the payload
	// return errmsg in the body as the error
	if rsp.StatusCode == 400 {
		errmsg, _ := ioutil.ReadAll(rsp.Body)
		return rsp, errors.New(string(errmsg))
	}

	// There was a problem with the connection
	if err != nil {
		fmt.Println("Error is ", err)
		return rsp, err
	}
	return rsp, nil
}

// GetTokens(login)
// Return a list of tokens for the current user
// or if login is specified use it
//
func (c *sonarcloudclient) GetTokens(name string) (*http.Response, error) {
	var url string
	if name != "" {
		options := "?"
		options += fmt.Sprintf(Login, name)
		url = c.URI + TokenSearch + options
	} else {
		url = c.URI + TokenSearch
	}

	req, err := http.NewRequest("GET", url, nil)
	resp, err := c.Client.Do(req)

	if err != nil {
		fmt.Println("Error is ", err)
		return resp, err
	}

	return resp, nil
}

// GetMetric(metric, branch string)
// Return an SVG badge for inclussion in HTML
//
// 	metric (required) is one of the following constants:
//    Bugs
//		CodeSmells
//		Coverage
//		DuplicatedLinesDensity
//		Ncloc
//		SqaleRating
//		AlertStatus
//		ReliabilityRating
//		SecurityRating
//		SqaleIndex
//
//  project  (required) project to produce bade for
//  branch (optional) a long living branch
//
func (c *sonarcloudclient) GetMetric(metric int, project, branch string) (*http.Response, error) {
	var url string
	options := "?"
	options += fmt.Sprintf(Metric, MetricName[metric])
	options += fmt.Sprintf("&"+Project, project)
	if branch != "" {
		options += fmt.Sprintf("&"+Branch, branch)
	}
	url = c.URI + BadgeMetric + options

	req, err := http.NewRequest("GET", url, nil)
	resp, err := c.Client.Do(req)

	if err != nil {
		fmt.Println("Error is ", err)
		return resp, err
	}

	return resp, nil
}

// GetQualityGate(project string) (*http.Response, error)
// Return an SVG badge for inclusion in HTML
// 	project (required) is a valid project name
//
func (c *sonarcloudclient) GetQualityGate(project string) (*http.Response, error) {
	var url string
	options := "?"
	options += fmt.Sprintf(Project, project)
	url = c.URI + QualityGate + options

	req, err := http.NewRequest("GET", url, nil)
	resp, err := c.Client.Do(req)

	if err != nil {
		fmt.Println("Error is ", err)
		return resp, err
	}

	return resp, nil
}
