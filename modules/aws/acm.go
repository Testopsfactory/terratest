package aws

import (
	"github.com/aws/aws-sdk-go/service/acm"
)

// Get the ACM certificate for the given domain name in the given region
func GetAcmCertificateArn(awsRegion string, certDomainName string) (string, error) {
	acmClient, err := NewAcmClient(awsRegion)
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
func NewAcmClient(awsRegion string) (*acm.ACM, error) {
	sess, err := GetAuthenticatedSession(awsRegion)
	if err != nil {
		return nil, err
	}

	return acm.New(sess), nil
}
