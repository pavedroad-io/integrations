// Package sonarcloud
package sonarcloud

// Constants used in URI and query parameters
const (
	// Default api server for SonarCloud
	DefaultScheme = "https://"
	// Default api server for SonarCloud
	DefaultHost = "sonarcloud.io"
	// Default api version
	DefaultAPI = "/api"

	// ProjectSearch URI
	ProjectSearch = DefaultAPI + "/projects/search"

	// ProjectCreate URI
	ProjectCreate = DefaultAPI + "/projects/create"

	// ProjectDelete URI
	ProjectDelete = DefaultAPI + "/projects/delete"

	// TokenSearch URI
	TokenSearch = DefaultAPI + "/user_tokens/search"

	// TokenCreate URI
	TokenCreate = DefaultAPI + "/user_tokens/generate"

	// TokenRevoke URI
	TokenRevoke = DefaultAPI + "/user_tokens/revoke"

	// BadgeMetric URI
	BadgeMetric = DefaultAPI + "/project_badges/measure"

	// QualityGate URI
	QualityGate = DefaultAPI + "/project_badges/quality_gate"

	// Query parameter strings
	Organization = "organization=%s"
	Project      = "project=%s"
	Projects     = "projects=%s"
	Metric       = "metric=%s"
	Name         = "name=%s"

	// KeyPrefix to append to SonarCloud Key to ensure uniqueness
	KeyPrefix = "PavedRoad_"

	// Testing constants
	projectKey  = "test123"
	projectName = "Test project 123"
	orgname     = "acme-demo"
	visibility  = "public"
)
