package fisalisgo_test

import (
	"context"
	"testing"

	fisalisgo "github.com/LillySchramm/FiSaLisGo"
)

func TestSearch(t *testing.T) {
	ctx := context.Background()

	res, err := fisalisgo.Search(ctx, "Vladimir Vladimirovich Putin")
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(res))
	}

	if res[0].Description != "President of the Russian Federation." {
		t.Fatalf("Expected 'President of the Russian Federation.', got '%s'", res[0].Description)
	}

	if res[0].Documents[0].Id != "2022/332 (OJ L53)" {
		t.Fatalf("Expected '2022/332 (OJ L53)', got '%s'", res[0].Documents[0].Id)
	}
}

func TestEmptySearch(t *testing.T) {
	ctx := context.Background()

	res, err := fisalisgo.Search(ctx, "Bernd das Brot")
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 0 {
		t.Fatalf("Expected 0 results, got %d", len(res))
	}
}
