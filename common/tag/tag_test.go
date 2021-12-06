package tag

// import (
// 	"testing"
// )

// func TestProcessOptions(t *testing.T) {

// 	tags := make(map[string]string)
// 	tags["tag1"] = "tag1-value"
// 	tags["tag2"] = "tag2-value"
// 	tags["tag3"] = "tag3-value"
// 	tags["tag4"] = "tag4-value"
// 	tags["tag5"] = "tag5-value"

// 	m := NewMatcher(tags)

// 	{
// 		tagsTest := make(map[string]string)
// 		tagsTest["tag1"] = "tag1-value"

// 		if !m.Match(tagsTest) {
// 			t.Fatalf("not expected")
// 		}
// 	}

// 	{
// 		tagsTest := make(map[string]string)
// 		tagsTest["tag1"] = "tag1-value"
// 		tagsTest["tag4"] = "tag4-value"

// 		if !m.Match(tagsTest) {
// 			t.Fatalf("not expected")
// 		}
// 	}

// 	{
// 		tagsTest := make(map[string]string)
// 		tagsTest["tag25"] = "tag25-value"

// 		if m.Match(tagsTest) {
// 			t.Fatalf("not expected")
// 		}
// 	}

// 	{
// 		tagsTest := make(map[string]string)
// 		tagsTest["tag25"] = "tag1-value"

// 		if m.Match(tagsTest) {
// 			t.Fatalf("not expected")
// 		}
// 	}

// 	{
// 		tagsTest := make(map[string]string)
// 		tagsTest["tag1"] = "tag1-BAD-value"

// 		if m.Match(tagsTest) {
// 			t.Fatalf("not expected")
// 		}
// 	}

// 	{
// 		tagsTest := make(map[string]string)
// 		tagsTest["tag1"] = "tag1-BAD-value"
// 		tagsTest["tag5"] = "tag5-value"

// 		if !m.Match(tagsTest) {
// 			t.Fatalf("not expected")
// 		}
// 	}

// }
