package command

import (
	"testing"
)

func standardJsonSelectTestData() []byte {
	return []byte(`{
  "users": [
    {
      "name": {
        "first": "Tom",
        "last": "Wright"
      },
      "flags": {
        "isBanned": false
      }
    },
    {
      "name": {
        "first": "Jim",
        "last": "Wright"
      },
      "flags": {
        "isBanned": true
      }
    },
    {
      "name": {
        "first": "Joe",
        "last": "Blogs"
      },
      "flags": {
        "isBanned": false
      }
    }
  ]
}`)
}

func TestSelectCommand(t *testing.T) {

	t.Run("TotalUsers", runTest(
		[]string{"-r", "json", "--pretty=false", "users.len()"},
		standardJsonSelectTestData(),
		newline([]byte(`3`)),
		nil,
		nil,
	))

	t.Run("TotalBannedUsers", runTest(
		[]string{"-r", "json", "--pretty=false", "users.all().filter(equal(flags.isBanned,true)).count()"},
		standardJsonSelectTestData(),
		newline([]byte(`1`)),
		nil,
		nil,
	))

	t.Run("TotalUnBannedUsers", runTest(
		[]string{"-r", "json", "--pretty=false", "users.all().filter(equal(flags.isBanned,false)).count()"},
		standardJsonSelectTestData(),
		newline([]byte(`2`)),
		nil,
		nil,
	))

	t.Run("NotBannedUsers", runTest(
		[]string{"-r", "json", "--pretty=false", "users.all().filter(equal(flags.isBanned,false)).name.first"},
		standardJsonSelectTestData(),
		newline([]byte(`"Tom"
"Joe"`)),
		nil,
		nil,
	))

	t.Run("BannedUsers", runTest(
		[]string{"-r", "json", "--pretty=false", "users.all().filter(equal(flags.isBanned,true)).name.first"},
		standardJsonSelectTestData(),
		newline([]byte(`"Jim"`)),
		nil,
		nil,
	))
}
