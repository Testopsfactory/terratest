package gcp

import "os"

var credsEnvVars = []string{
	"GOOGLE_CREDENTIALS",
	"GOOGLE_CLOUD_KEYFILE_JSON",
	"GCLOUD_KEYFILE_JSON",
	"GOOGLE_USE_DEFAULT_CREDENTIALS",
}

var projectEnvVars = []string{
	"GOOGLE_PROJECT",
	"GOOGLE_CLOUD_PROJECT",
	"GOOGLE_CLOUD_PROJECT_ID",
	"GCLOUD_PROJECT",
	"CLOUDSDK_CORE_PROJECT",
}

var regionEnvVars = []string{
	"GOOGLE_REGION",
	"GCLOUD_REGION",
	"CLOUDSDK_COMPUTE_REGION",
}

// GetGoogleCredentialsFromEnvVar returns the Credentials for use with testing.
func GetGoogleCredentialsFromEnvVar() string {
	return multiEnvSearch(credsEnvVars)
}

// GetGoogleProjectIDFromEnvVar returns the Project Id for use with testing.
func GetGoogleProjectIDFromEnvVar() string {
	return multiEnvSearch(projectEnvVars)
}

// GetGoogleRegionFromEnvVar returns the Region for use with testing.
func GetGoogleRegionFromEnvVar() string {
	return multiEnvSearch(regionEnvVars)
}

func multiEnvSearch(ks []string) string {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}
