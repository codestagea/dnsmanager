package diff

import (
	"fmt"
	"testing"
	"time"
)

type DiffTestS struct {
	Name      string `diff:"名字"`
	Value     string `diff:"值"`
	Ttl       int    `diff:"ttl"`
	CreatedAt time.Time
}

func TestDiffStruct(t *testing.T) {
	before := DiffTestS{
		Name:      "diff_test",
		Value:     "origin",
		Ttl:       10,
		CreatedAt: time.Now(),
	}

	after := DiffTestS{
		Name:      "diff_test",
		Value:     "newValue",
		Ttl:       20,
		CreatedAt: time.Now(),
	}

	results := DiffStruct(before, after)
	for _, r := range results {
		fmt.Println(r.String())
	}
	if len(results) != 2 {
		t.Errorf("expect 2, actual %d", len(results))
		return
	}
	if results[0].String() != "值: origin -> newValue" {
		t.Errorf("value 0 fail: %s", results[0].String())
		return
	}
	if results[1].String() != "ttl: 10 -> 20" {
		t.Errorf("value 1 fail: %s", results[1].String())
		return
	}
}

func TestDiffStructPtr(t *testing.T) {
	before := DiffTestS{
		Name:      "diff_test",
		Value:     "origin",
		Ttl:       10,
		CreatedAt: time.Now(),
	}

	after := DiffTestS{
		Name:      "diff_test",
		Value:     "newValue",
		Ttl:       20,
		CreatedAt: time.Now(),
	}

	results := DiffStruct(&before, &after)
	for _, r := range results {
		fmt.Println(r.String())
	}

	if len(results) != 2 {
		t.Errorf("expect 2, actual %d", len(results))
		return
	}
	if results[0].String() != "值: origin -> newValue" {
		t.Errorf("value 0 fail: %s", results[0].String())
		return
	}
	if results[1].String() != "ttl: 10 -> 20" {
		t.Errorf("value 1 fail: %s", results[1].String())
		return
	}
}
