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

	t.Run("TotalUsersLen", runTest(
		[]string{"-r", "json", "--pretty=false", "users.len()"},
		standardJsonSelectTestData(),
		newline([]byte(`3`)),
		nil,
		nil,
	))

	t.Run("TotalUsersCount", runTest(
		[]string{"-r", "json", "--pretty=false", "users.all().count()"},
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

	t.Run("TotalNotBannedUsers", runTest(
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

	t.Run("Issue258", runTest(
		[]string{"-r", "json", "--pretty=false", "-w", "csv", "phones.all().mapOf(make,make,model,model,first,parent().parent().user.name.first,last,parent().parent().user.name.last).merge()"},
		[]byte(`{
	  "id": "1234",
	  "user": {
	    "name": {
	      "first": "Tom",
	      "last": "Wright"
	    }
	  },
	  "favouriteNumbers": [
	    1, 2, 3, 4
	  ],
	  "favouriteColours": [
	    "red", "green"
	  ],
	  "phones": [
	    {
	      "make": "OnePlus",
	      "model": "8 Pro"
	    },
	    {
	      "make": "Apple",
	      "model": "iPhone 12"
	    }
	  ]
	}`),
		newline([]byte(`first,last,make,model
Tom,Wright,OnePlus,8 Pro
Tom,Wright,Apple,iPhone 12`)),
		nil,
		nil,
	))
}
