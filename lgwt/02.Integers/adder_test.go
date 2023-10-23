package integers

import (
	"fmt"
	"testing"
)

func TestAdder(t *testing.T) {
	sum := Add(2, 2)
	expected := 4
	if sum != expected {
		t.Errorf("\n"+
			"expected:\n"+
			"'%d'\n"+
			"actual:\n"+
			"'%d'\n",
			expected, sum)
	}
}

// Пример запустится в составе тестов, если будет присутствовать
// комментарий "// Output: 6"
func ExampleAdd() {
	sum := Add(1, 5)
	fmt.Println(sum)
	// Output: 6
}
