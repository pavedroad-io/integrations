# SonarCloud Integration

The PavedRoad SonarCloud integration supports the following capabilities;

Search, Create and Delete projects
Search, Create and Revoke tokens
Generate SVG Badges for metrics
Generate SVG Badge for quality gate

# Usage

## SONARCLOUD_TOKEN
Set the SONARCLOUD_TOKEN environment variable with a user token that has administrative privileges

## Create a client
The New methods set standard defaults and also creates an HTTP client.

```go
  testClient = sonarcloudclient{}
  err := testClient.New(token)
  if err != nil {
    t.Errorf("Expected err to be nil Got %v\n", err)
  }
```
## Projects

```go
  // Get a project
  rsp, err := testClient.GetProject(organaization, projectKey)

  var prj ProjectSearchResponse
  err = json.Unmarshal(project, &prj)

  if err != nil {
   t.Errorf("Unmarshal failed Got %v\n", err)
  }

  // Create new project
  p := NewProject{
    Organization: orgname,
    Name:         projectName,
    Project:      projectKey,
    Visibility:   visibility,
  }


  rsp, err := testClient.CreateProject(p)
  if err != nil {
    t.Errorf("Expected err to be nil Got %s\n", err)
    return
  }

  // Delete a project
  CreateProjectIfItDoesntExists(t)
  rsp, err := testClient.DeleteProject(projectKey)
  if err != nil {
    t.Errorf("Expected err to be nil Got %v\n", err)
  }
```
## Tokens


```go
  // Create a token

  rsp, err := testClient.CreateToken(tokenName)
  if err != nil {
    t.Errorf("Expected err to be nil Got %s\n", err)
    return
  }

  var tk NewTokenResponse
  err = json.Unmarshal(token, &tk)

  if err != nil {
    t.Errorf("Unmarshal failed Got %v\n", err)
  }

  // Get a token, loginName is optional and defaults to current user
  loginName := ""
  rsp, err := testClient.GetTokens(loginName)
  if err != nil {
    t.Errorf("Expected err to be nil Got %v\n", err)
  }

  var tk GetTokenResponse
  err = json.Unmarshal(tkList, &tk)

  // Revoke a token
  rsp, err := testClient.RevokeToken(tokenName)
  if err != nil {
    t.Errorf("Expected err to be nil Got %v\n", err)
  }
```

## Badges
Returns SVG snip for including in HTML

Valid metrics are:
- Bugs,
-     CodeSmells,
-     Coverage,
-     DuplicatedLinesDensity,
-     Ncloc,
-     SqaleRating,
-     AlertStatus,
-     ReliabilityRating,
-     SecurityRating,
-     SqaleIndex
-     Vulnerabilities,

```go

  //Metric badges
  rsp, err := testClient.GetMetric(m, projectKey, branch)
  if err != nil {
    t.Errorf("Expected err to be nil Got %v\n", err)
  }

  svg, err := ioutil.ReadAll(rsp.Body)
  if err != nil {
    t.Errorf("Expected err to be nil Got %v\n", err)
  }
  svgStr := string(svg)

  // Quality Gate Badge
  rsp, err := testClient.GetQualityGate(projectKey)
  if err != nil {
    t.Errorf("Expected err to be nil Got %v\n", err)
  }

  svg, err := ioutil.ReadAll(rsp.Body)
  if err != nil {
    t.Errorf("Expected err to be nil Got %v\n", err)
  }
  svgStr := string(svg)
```


