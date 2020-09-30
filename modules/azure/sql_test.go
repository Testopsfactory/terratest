// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.
package azure

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/require"
)

func TestGetSQLServerIDE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	serverName := ""
	subscriptionID := ""

	_, err := azure.GetSQLServerIDE(t, resGroupName, serverName, subscriptionID)
	require.Error(t, err)
}

func TestGetSQLServerFullDomainNameE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	serverName := ""
	subscriptionID := ""

	_, err := azure.GetSQLServerFullDomainNameE(t, resGroupName, serverName, subscriptionID)
	require.Error(t, err)
}

func TestGetSQLServerStateE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	serverName := ""
	subscriptionID := ""

	_, err := azure.GetSQLServerStateE(t, resGroupName, serverName, subscriptionID)
	require.Error(t, err)
}

func TestGetDatabaseIDE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	serverName := ""
	dbName := ""
	subscriptionID := ""

	_, err := azure.GetDatabaseIDE(t, resGroupName, serverName, dbName, subscriptionID)
	require.Error(t, err)
}

func TestGetDatabaseStatusE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	serverName := ""
	dbName := ""
	subscriptionID := ""

	_, err := azure.GetDatabaseStatusE(t, resGroupName, serverName, dbName, subscriptionID)
	require.Error(t, err)
}
