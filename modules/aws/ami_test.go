package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUbuntu1404AmiReturnsSomeAmi(t *testing.T) {
	t.Parallel()

	amiId := GetUbuntu1404Ami(t, "us-east-1")
	assert.Regexp(t, "^ami-[[:alnum:]]+$", amiId)
}

func TestGetUbuntu1604AmiReturnsSomeAmi(t *testing.T) {
	t.Parallel()

	amiId := GetUbuntu1604Ami(t, "us-west-1")
	assert.Regexp(t, "^ami-[[:alnum:]]+$", amiId)
}

func TestGetCentos7AmiReturnsSomeAmi(t *testing.T) {
	t.Parallel()

	amiId := GetCentos7Ami(t, "eu-west-1")
	assert.Regexp(t, "^ami-[[:alnum:]]+$", amiId)
}

func TestGetAmazonLinuxAmiReturnsSomeAmi(t *testing.T) {
	t.Parallel()

	amiId := GetAmazonLinuxAmi(t, "ap-southeast-1")
	assert.Regexp(t, "^ami-[[:alnum:]]+$", amiId)
}

func TestGetEcsOptimizedAmazonLinuxAmiEReturnsSomeAmi(t *testing.T) {
	t.Parallel()

	amiId := GetEcsOptimizedAmazonLinuxAmi(t, "us-east-2")
	assert.Regexp(t, "^ami-[[:alnum:]]+$", amiId)
}
