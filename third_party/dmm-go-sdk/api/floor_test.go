package api

import (
	"reflect"
	"testing"
)

func TestNewFloorService(t *testing.T) {
	affiliateID := "foobar-999"
	apiID := "TXpEZ5D4T2xB3J5cuSLf"

	srv := NewFloorService(affiliateID, apiID)
	if srv.AffiliateID != affiliateID {
		t.Fatalf("FloorService.AffiliateID is expected to equal the input value(affiliateID)")
	}

	if srv.ApiID != apiID {
		t.Fatalf("FloorService.ApiID is expected to equal the input value(apiID)")
	}
}

func TestBuildRequestURLInFloorService(t *testing.T) {
	var srv *FloorService
	var u string
	var err error
	var expected string

	srv = dummyFloorService()
	u, err = srv.BuildRequestURL()
	expected = APIBaseURL + "/FloorList?affiliate_id=" + DummyAffliateID + "&api_id=" + DummyAPIID
	if u != expected {
		t.Fatalf("FloorService.BuildRequestURL is expected to equal the expected value.\nexpected:%s\nactual:  %s", expected, u)
	}
	if err != nil {
		t.Fatalf("FloorService.BuildRequestURL is not expected to have error")
	}
}

func TestBuildRequestURLWithoutApiIDInInFloorService(t *testing.T) {
	srv := dummyFloorService()
	srv.ApiID = ""
	u, err := srv.BuildRequestURL()
	if u != "" {
		t.Fatalf("FloorService.BuildRequestURL is expected empty if API ID is not set.")
	}
	if err == nil {
		t.Fatalf("FloorService.BuildRequestURL is expected to return error.")
	}
}

func TestBuildRequestURLWithWrongAffiliateIDInFloorService(t *testing.T) {
	srv := dummyFloorService()
	srv.AffiliateID = "fizzbizz-100"
	u, err := srv.BuildRequestURL()
	if u != "" {
		t.Fatalf("FloorService.BuildRequestURL is expected empty if wrong Affiliate ID is set.")
	}
	if err == nil {
		t.Fatalf("FloorService.BuildRequestURL is expected to return error.")
	}
}

func TestExcuteRequestFloorAPIToServer(t *testing.T) {
	if !RequestAvailable {
		t.Skip("Not set valid credentials")
	}

	srv := NewFloorService(TestAffiliateID, TestAPIID)
	rst, err := srv.Execute()
	if err != nil {
		t.Skip("Maybe, The network is down.")
	}

	if reflect.TypeOf(rst).String() != "*api.FloorResponse" {
		t.Fatalf("FloorService.Execute is expected to return *api.FloorResponse")
	}
}

func dummyFloorService() *FloorService {
	return NewFloorService(DummyAffliateID, DummyAPIID)
}
