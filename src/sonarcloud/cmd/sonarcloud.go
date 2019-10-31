// Package sonarcloud
package sonarcloud

import (
	"crypto/tls"
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
	size int `json:"pageSize"`

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

// NewProject
//   Used to create a new project
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

// NewProjectResponse
//   Used to create a new project
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

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	c.Client = &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
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
