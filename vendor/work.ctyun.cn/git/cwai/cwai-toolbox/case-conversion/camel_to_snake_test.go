package case_conversion

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testCase struct {
	CamelCase         string
	ExpectedSnakeCase string
}

type camelTestCases []testCase

func (cases camelTestCases) Test(t *testing.T) {
	for _, testCase := range cases {
		assert.Equal(t, testCase.ExpectedSnakeCase, CamelToSnake(testCase.CamelCase))
	}
}

func TestIdentities(t *testing.T) {
	camelTestCases{
		{"username", "username"},
		{"parameter1", "parameter1"},
	}.Test(t)
}

func TestConsequentLargeCases(t *testing.T) {
	camelTestCases{
		{"ID", "id"},
		{"eventID", "event_id"},
		{"NSString", "ns_string"},
	}.Test(t)
}

func TestSmallCamelCases(t *testing.T) {
	camelTestCases{
		{"dataPipelineID", "data_pipeline_id"},
		{"someIDIsOkayHere", "some_id_is_okay_here"},
		{"mixedIDIsAlsoOkay", "mixed_id_is_also_okay"},
	}.Test(t)
}

func TestBigCamelCases(t *testing.T) {
	camelTestCases{
		{"HelloIndianMIFans", "hello_indian_mi_fans"},
		{"DoYouLikeMiBand", "do_you_like_mi_band"},
	}.Test(t)
}
