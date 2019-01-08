package parser

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"

	cfg "github.com/mbrt/gmailctl/pkg/config/v1alpha2"
)

func readConfig(t *testing.T, path string) cfg.Config {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var res cfg.Config
	if err := yaml.UnmarshalStrict(b, &res); err != nil {
		t.Fatal(err)
	}
	return res
}

func TestParse(t *testing.T) {
	conf := readConfig(t, "testdata/example.yaml")
	expected := []Rule{
		{
			Criteria: and(
				fn(FunctionList, OperationOr,
					"list1",
					"list2",
					"list3",
				),
				not(
					fn(FunctionTo, OperationOr,
						"pippo@gmail.com",
						"pippo@hotmail.com",
					),
				),
			),
			Actions: cfg.Actions{
				Labels:  []string{"maillist"},
				Archive: true,
			},
		},
		{
			Criteria: and(
				fn(FunctionTo, OperationAnd, "myalias@gmail.com"),
				fn(FunctionList, OperationOr,
					"list1",
					"list2",
					"list3",
				),
			),
			Actions: cfg.Actions{MarkImportant: true},
		},
		{
			Criteria: fn(FunctionFrom, OperationOr,
				"spammer1", "spammer2",
			),
			Actions: cfg.Actions{Delete: true},
		},
		{
			Criteria: fn(FunctionTo, OperationOr,
				"pippo+spammy@gmail.com",
			),
			Actions: cfg.Actions{Delete: true},
		},
		{
			Criteria: fn(FunctionSubject, OperationOr,
				"spam mail",
			),
			Actions: cfg.Actions{Delete: true},
		},
		{
			Criteria: fn(FunctionHas, OperationOr,
				"buy this thing",
				"very important!!!",
			),
			Actions: cfg.Actions{Delete: true},
		},
	}
	got, err := Parse(conf)
	assert.Nil(t, err)
	assert.Equal(t, expected, got)
}
