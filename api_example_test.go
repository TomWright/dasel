package dasel_test

import (
	"context"
	"fmt"
	"github.com/tomwright/dasel/v3"
	"github.com/tomwright/dasel/v3/execution"
)

func ExampleSelect() {
	myData := map[string]any{
		"users": []map[string]any{
			{"name": "Alice", "age": 30},
			{"name": "Bob", "age": 25},
			{"name": "Tom", "age": 40},
		},
	}
	query := `users.filter(age > 27).map(name)...`
	selectResult, numResults, err := dasel.Select(context.Background(), myData, query, execution.WithUnstable())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d results:\n", numResults)

	// You should validate the type assertion in real code.
	selectResults := selectResult.([]any)

	// Results can be of various types, handle accordingly.
	for _, result := range selectResults {
		fmt.Println(result)
	}

	// Output:
	// Found 2 results:
	// Alice
	// Tom
}
