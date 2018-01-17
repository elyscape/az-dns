package helpers

import (
	"fmt"
	"testing"
)

type recordNameTestCase struct {
	hostname       string
	zone           string
	relative       bool
	expectedRecord string
}

var recordNameTests = []recordNameTestCase{
	// Apex
	{"", "example.com", false, "@"},
	{"@", "example.com", false, "@"},
	{"", "example.com", true, "@"},
	{"@", "example.com", true, "@"},
	{"example.com", "example.com", false, "@"},
	{"example.com.", "example.com", false, "@"},

	// Subdomain
	{"sub.example.com", "example.com", false, "sub"},
	{"sub.example.com.", "example.com", false, "sub"},
	{"sub", "example.com", false, "sub"},
	{"sub.", "example.com", false, "sub"},

	// Relative subdomain
	{"example.com", "example.com", true, "example.com"},
	{"example.com.", "example.com", true, "example.com"},
	{"sub.example.com", "example.com", true, "sub.example.com"},
	{"sub.example.com.", "example.com", true, "sub.example.com"},
	{"sub", "example.com", true, "sub"},
	{"sub.", "example.com", true, "sub"},
}

func TestGenerateRecordName(t *testing.T) {
	for _, testCase := range recordNameTests {
		var name string
		if !testCase.relative {
			name = fmt.Sprintf("%v@%v", testCase.hostname, testCase.zone)
		} else {
			name = fmt.Sprintf("%v(rel)@%v", testCase.hostname, testCase.zone)
		}
		t.Run(name, func(t *testing.T) { testGenerateRecordName(t, testCase) })
	}
}

func testGenerateRecordName(t *testing.T, testCase recordNameTestCase) {
	result := GenerateRecordName(testCase.hostname, testCase.zone, testCase.relative)
	if result != testCase.expectedRecord {
		t.Errorf(`got "%v", expected "%v"`, result, testCase.expectedRecord)
	}
}
