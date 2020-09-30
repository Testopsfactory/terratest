// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.
package azure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when CRUD methods are introduced for Azure SQL DB, these tests can be extended
*/

func TestGetSQLServerIDE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	serverName := ""
	subscriptionID := ""

	_, err := GetSQLServerIDE(t, resGroupName, serverName, subscriptionID)
	require.Error(t, err)
}

func TestGetSQLServerFullDomainNameE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	serverName := ""
	subscriptionID := ""

	_, err := GetSQLServerFullDomainNameE(t, resGroupName, serverName, subscriptionID)
	require.Error(t, err)
}

func TestGetSQLServerStateE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	serverName := ""
	subscriptionID := ""

	_, err := GetSQLServerStateE(t, resGroupName, serverName, subscriptionID)
	require.Error(t, err)
}

func TestGetDatabaseIDE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	serverName := ""
	dbName := ""
	subscriptionID := ""

	_, err := GetDatabaseIDE(t, resGroupName, serverName, dbName, subscriptionID)
	require.Error(t, err)
}

func TestGetDatabaseStatusE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	serverName := ""
	dbName := ""
	subscriptionID := ""

	_, err := GetDatabaseStatusE(t, resGroupName, serverName, dbName, subscriptionID)
	require.Error(t, err)
}
