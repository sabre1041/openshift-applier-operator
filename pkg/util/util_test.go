package util

import "testing"

func TestParsing(t *testing.T) {
	tables := []struct {
		path      string
		namespace string
		name      string
		token     string
	}{
		{"/testnamespace/testresource/mytoken", "testnamespace", "testresource", "mytoken"},
	}

	for _, table := range tables {
		namespace, name, token, err := ParseQueryString(table.path)

		if err != nil {

		}

		// Verify no errors occur
		if err != nil {
			t.Errorf("Error occurred processing '%s': %v", table.path, err)
		}

		// Validate Namespace
		if namespace != table.namespace {
			t.Errorf("Namespace from path '%s' was incorrect, got: %s, want: %s.", table.path, namespace, table.namespace)
		}

		// Validate Name
		if namespace != table.namespace {
			t.Errorf("Resource Name from path '%s' was incorrect, got: %s, want: %s.", table.path, name, table.name)
		}

		// Validate Token
		if token != table.token {
			t.Errorf("Token from path '%s' was incorrect, got: %s, want: %s.", table.path, token, table.token)
		}
	}
}
