package azure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when CRUD methods are introduced for Azure Synapse, these tests can be extended
*/

func TestGetDataFactoryE(t *testing.T) {
	t.Parallel()

	resGroupName := "terratest-datafactory-resource"
	subscriptionID := "00fb78cc-7201-4e1c-8203-2b2e1390309a"
	dataFactoryName := "datafactoryresource"

	_, err := GetDataFactoryE(t, subscriptionID, resGroupName, dataFactoryName)
	require.Error(t, err)
}
