package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
)

func CreateAcmClient(awsRegion string) (*acm.ACM, error) {
	awsConfig, err := CreateAwsConfig(awsRegion)
	if err != nil {
		return nil, err
	}

	return acm.New(session.New(), awsConfig), nil
}

// This exists because it is not possible to fully automate the requesting and approval of acm certificates. Therefore,
// we have created a *.gruntwork.io certificate in each region and approved it through DNS verification. These certs
// can now be used in automated tests that require an acm certificate.
func GetAcmCertificateArn(awsRegion string, certDomainName string) (string, error) {
	acmClient, err := CreateAcmClient(awsRegion)
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
