package sonarcloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
)

var testClient sonarcloudclient

func TestMain(t *testing.T) {
	var token string
	// Get token so we can run tests
	envVar := os.Getenv("SONARCLOUD_TOKEN")
	if envVar != "" {
		token = envVar
	} else {
		log.Println("Need SONARCLOUD_TOKEN set to run tests")
		os.Exit(-1)
	}

	// Setup the client
	testClient = sonarcloudclient{}
	testClient.New(token)
}

// TestGetProject make sure we can retrive a project
func TestGetProject(t *testing.T) {
	CreateProjectIfItDoesntExists(t)
	x := projectKey
	rsp, err := testClient.GetProject(orgname, x)
	if err != nil {
		t.Errorf("Expected err to be nil Got %v\n", err)
	}

	checkResponseCode(t, http.StatusOK, rsp.StatusCode)

	project, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.Errorf("Expected err to be nil Got %v\n", err)
	}

	var prj ProjectSearchResponse
	err = json.Unmarshal(project, &prj)

	if err != nil {
		t.Errorf("Unmarshal failed Got %v\n", err)
	}

	if len(prj.Components) > 0 && prj.Components[0].Key != x {
		t.Errorf("Expected key to be '"+x+"'. Got '%v'", prj.Components[0].Key)
	}
}

func TestCreateProject(t *testing.T) {
	p := NewProject{
		Organization: orgname,
		Name:         projectName,
		Project:      projectKey,
		Visibility:   visibility,
	}

	DeleteProjectIfExists(projectKey)

	rsp, err := testClient.CreateProject(p)
	if err != nil {
		t.Errorf("Expected err to be nil Got %s\n", err)
		return
	}

	// Note: SonarCloud returns the wrong status code
	checkResponseCode(t, http.StatusOK, rsp.StatusCode)

	if http.StatusOK != rsp.StatusCode {
		errrsp, _ := ioutil.ReadAll(rsp.Body)
		fmt.Println(string(errrsp))
	}

	project, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.Errorf("Expected err to be nil Got %v\n", err)
	}

	var prj NewProjectResponse
	err = json.Unmarshal(project, &prj)

	if err != nil {
		t.Errorf("Unmarshal failed Got %v\n", err)
	}

	e := KeyPrefix + projectKey
	if prj.Project.Key != e {
		t.Errorf("Expected key to be "+e+". Got '%v'", prj.Project.Key)
	}
}

// TestDeleteProject make sure we can delete a project
func TestDeleteProject(t *testing.T) {
	CreateProjectIfItDoesntExists(t)
	rsp, err := testClient.DeleteProject(projectKey)
	if err != nil {
		t.Errorf("Expected err to be nil Got %v\n", err)
	}

	checkResponseCode(t, http.StatusNoContent, rsp.StatusCode)
}

// DeleteProjectIfExists where project name is p
func DeleteProjectIfExists(p string) {
	_, err := testClient.DeleteProject(p)

	if err != nil {
		fmt.Println("DeleteProjectIfExists ", p, err)
	}
	return
}

// CreateProjectIfItDoesntExists where project name is p
// Use project defaults
func CreateProjectIfItDoesntExists(t *testing.T) {
	p := NewProject{
		Organization: orgname,
		Name:         projectName,
		Project:      projectKey,
		Visibility:   visibility,
	}
	_, err := testClient.CreateProject(p)

	if err != nil {
		// Ignore if it is because the records already exists
		if strings.Contains(err.Error(), "key already exists") {
			return
		}

		t.Errorf("CreateProjectIfItDoesntExists %v %v\n", p, err.Error())
	}
	return
}

// Check the expected and actual response codes
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
