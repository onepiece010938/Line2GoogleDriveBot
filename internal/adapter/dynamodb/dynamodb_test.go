package dynamodb

import (
	"os"
	"testing"
)

var testTableBasics *TableBasics

func TestMain(m *testing.M) {
	// google_oauth
	testTableBasics = NewTableBasics("google-oauth")
	// change to local dynamodb
	testTableBasics.DynamoDbClient = CreateLocalClient(8000)
	os.Exit(m.Run())
}

func TestTableExists(t *testing.T) {
	exits, err := testTableBasics.TableExists()
	t.Log("Exists:", exits)
	t.Log("ERROR:", err)
}

func TestListTables(t *testing.T) {
	tableNames, err := testTableBasics.ListTables()
	t.Log("tableNames:", tableNames)
	t.Log("ERROR:", err)

}
