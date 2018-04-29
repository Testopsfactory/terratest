package aws

import (
	"github.com/aws/aws-sdk-go/service/acm"
	"testing"
)

// Get the ACM certificate for the given domain name in the given region
func GetAcmCertificateArn(t *testing.T, awsRegion string, certDomainName string) string {
	arn, err := GetAcmCertificateArnE(t, awsRegion, certDomainName)
	if err != nil {
		t.Fatal(err)
	}
	return arn
}

// Get the ACM certificate for the given domain name in the given region
func GetAcmCertificateArnE(t *testing.T, awsRegion string, certDomainName string) (string, error) {
	acmClient, err := NewAcmClientE(t, awsRegion)
	if err != nil {
		return "", err
	}

	result, err := acmClient.ListCertificates(&acm.ListCertificatesInput{})
	if err != nil {
		return "", err
	}

	for _, summary := range result.CertificateSummaryList {
		if *summary.DomainName == certDomainName {
			return *summary.CertificateArn, nil
		}
	}

	return "", nil
}

// Create a new ACM client
func NewAcmClient(t *testing.T, region string) *acm.ACM {
	client, err := NewAcmClientE(t, region)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// Create a new ACM client
func NewAcmClientE(t *testing.T, awsRegion string) (*acm.ACM, error) {
	sess, err := NewAuthenticatedSession(awsRegion)
	if err != nil {
		return nil, err
	}

	return acm.New(sess), nil
}
