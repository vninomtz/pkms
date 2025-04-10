package index

import (
	"fmt"
	"testing"
)

func TestTokenize(t *testing.T) {
	input := "Hello word! con text-full it's and 928376 cases && for piñas"
	expected := 10

	res := Tokenize(input)
	fmt.Printf("Result %d: \n", len(res))
	for i, t := range res {
		if i == (len(res) - 1) {
			fmt.Printf("%s\n", t.Value)
		} else {
			fmt.Printf("%s,", t.Value)
		}
	}
	if len(res) != expected {
		t.Errorf("Expected %d, got %d", expected, len(res))
	}
}
