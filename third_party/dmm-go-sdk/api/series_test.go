package api

import (
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

func TestNewSeriesService(t *testing.T) {
	affiliateID := DummyAffliateID
	apiID := DummyAPIID

	srv := NewSeriesService(affiliateID, apiID)
	if srv.AffiliateID != affiliateID {
		t.Fatalf("SeriesService.AffiliateID is expected to equal the input value(affiliateID)")
	}

	if srv.ApiID != apiID {
		t.Fatalf("SeriesService.ApiID is expected to equal the input value(apiID)")
	}
}

func TestSetLengthInSeriesService(t *testing.T) {
	srv := dummySeriesService()
	var length int64 = 10
	srv.SetLength(length)

	if srv.Length != length {
		t.Fatalf("SeriesService.Length is expected to equal the input value(length)")
	}
}

func TestSetHitsInSeriesService(t *testing.T) {
	srv := dummySeriesService()
	var hits int64 = 10
	srv.SetHits(hits)

	if srv.Length != hits {
		t.Fatalf("SeriesService.Length is expected to equal the input value(hits)")
	}
}

func TestSetOffsetInSeriesService(t *testing.T) {
	srv := dummySeriesService()
	var offset int64 = 10
	srv.SetOffset(offset)

	if srv.Offset != offset {
		t.Fatalf("SeriesService.Offset is expected to equal the input value(offset)")
	}
}

func TestValidateLengthInSeriesService(t *testing.T) {
	srv := dummySeriesService()

	var target int64

	target = 1
	srv.SetLength(target)
	if srv.ValidateLength() == false {
		t.Fatalf("SeriesService.ValidateLength is expected TRUE.")
	}

	target = DefaultAPILength
	srv.SetLength(target)
	if srv.ValidateLength() == false {
		t.Fatalf("SeriesService.ValidateLength is expected TRUE.")
	}

	target = DefaultMaxLength
	srv.SetLength(target)
	if srv.ValidateLength() == false {
		t.Fatalf("SeriesService.ValidateLength is expected TRUE.")
	}

	target = DefaultMaxLength + 1
	srv.SetLength(target)
	if srv.ValidateLength() == true {
		t.Fatalf("SeriesService.ValidateLength is expected FALSE.")
	}

	target = 0
	srv.SetLength(target)
	if srv.ValidateLength() == true {
		t.Fatalf("SeriesService.ValidateLength is expected FALSE.")
	}

	target = -1
	srv.SetLength(target)
	if srv.ValidateLength() == true {
		t.Fatalf("SeriesService.ValidateLength is expected FALSE.")
	}
}

func TestValidateOffsetInSeriesService(t *testing.T) {
	srv := dummySeriesService()

	var target int64

	target = 1
	srv.SetOffset(target)
	if srv.ValidateOffset() == false {
		t.Fatalf("SeriesService.ValidateOffset is expected TRUE.")
	}

	target = 100
	srv.SetOffset(target)
	if srv.ValidateOffset() == false {
		t.Fatalf("SeriesService.ValidateOffset is expected TRUE.")
	}

	target = 0
	srv.SetOffset(target)
	if srv.ValidateOffset() == true {
		t.Fatalf("SeriesService.ValidateOffset is expected FALSE.")
	}

	target = -1
	srv.SetOffset(target)
	if srv.ValidateOffset() == true {
		t.Fatalf("SeriesService.ValidateOffset is expected FALSE.")
	}
}

func TestBuildRequestURLInSeriesService(t *testing.T) {
	var srv *SeriesService
	var u string
	var err error
	var expected string

	srv = dummySeriesService()
	srv.SetFloorID("40")
	u, err = srv.BuildRequestURL()
	expected = APIBaseURL + "/SeriesSearch?affiliate_id=" + DummyAffliateID + "&api_id=" + DummyAPIID + "&floor_id=40" + "&hits=" + strconv.FormatInt(DefaultAPILength, 10) + "&offset=" + strconv.FormatInt(DefaultAPIOffset, 10)
	if u != expected {
		t.Fatalf("SeriesService.BuildRequestURL is expected to equal the expected value.\nexpected:%s\nactual:  %s", expected, u)
	}
	if err != nil {
		t.Fatalf("SeriesService.BuildRequestURL is not expected to have error")
	}

	srv = dummySeriesService()
	srv.SetLength(0)
	srv.SetOffset(0)
	srv.SetFloorID("40")
	u, err = srv.BuildRequestURL()
	expected = APIBaseURL + "/SeriesSearch?affiliate_id=" + DummyAffliateID + "&api_id=" + DummyAPIID + "&floor_id=40"
	expectedBase := expected
	if u != expected {
		t.Fatalf("SeriesService.BuildRequestURL is expected to equal the expected value.\nexpected:%s\nactual:  %s", expected, u)
	}
	if err != nil {
		t.Fatalf("SeriesService.BuildRequestURL is not expected to have error")
	}

	srv.SetLength(-1)
	u, err = srv.BuildRequestURL()
	if u != "" {
		t.Fatalf("SeriesService.BuildRequestURL is expected empty if error occurs.")
	}
	if err == nil {
		t.Fatalf("SeriesService.BuildRequestURL is expected to return error.")
	}
	srv.SetLength(0)

	srv.SetOffset(-1)
	u, err = srv.BuildRequestURL()
	if u != "" {
		t.Fatalf("SeriesService.BuildRequestURL is expected empty if error occurs.")
	}
	if err == nil {
		t.Fatalf("SeriesService.BuildRequestURL is expected to return error.")
	}
	srv.SetOffset(0)

	srv.SetFloorID("")
	u, err = srv.BuildRequestURL()
	if u != "" {
		t.Fatalf("SeriesService.BuildRequestURL is expected empty if error occurs.")
	}
	if err == nil {
		t.Fatalf("SeriesService.BuildRequestURL is expected to return error.")
	}
	srv.SetFloorID("40")

	srv.SetInitial("あ")
	expected = expectedBase + "&initial=" + url.QueryEscape("あ")
	u, err = srv.BuildRequestURL()
	if u != expected {
		t.Fatalf("SeriesService.BuildRequestURL is expected to equal the expected value.\nexpected:%s\nactual:  %s", expected, u)
	}
	if err != nil {
		t.Fatalf("SeriesService.BuildRequestURL is not expected to have error")
	}
	srv.SetInitial("")
}

func TestBuildRequestURLWithoutApiIDInSeriesService(t *testing.T) {
	srv := dummySeriesService()
	srv.ApiID = ""
	u, err := srv.BuildRequestURL()
	if u != "" {
		t.Fatalf("SeriesService.BuildRequestURL is expected empty if API ID is not set.")
	}
	if err == nil {
		t.Fatalf("SeriesService.BuildRequestURL is expected to return error.")
	}
}

func TestBuildRequestURLWithWrongAffiliateIDInSeriesService(t *testing.T) {
	srv := dummySeriesService()
	srv.AffiliateID = "fizzbizz-100"
	u, err := srv.BuildRequestURL()
	if u != "" {
		t.Fatalf("SeriesService.BuildRequestURL is expected empty if wrong Affiliate ID is set.")
	}
	if err == nil {
		t.Fatalf("SeriesService.BuildRequestURL is expected to return error.")
	}
}

func TestExcuteRequestSeriesAPIToServer(t *testing.T) {
	if !RequestAvailable {
		t.Skip("Not set valid credentials")
	}

	srv := NewSeriesService(TestAffiliateID, TestAPIID)
	srv.SetFloorID("40")
	srv.SetInitial("あ")
	srv.SetLength(100)
	srv.SetOffset(1)

	rst, err := srv.Execute()
	if err != nil {
		t.Skip("Maybe, The network is down.")
	}

	if reflect.TypeOf(rst).String() != "*api.SeriesResponse" {
		t.Fatalf("SeriesService.Execute is expected to return *api.SeriesResponse")
	}
}

func dummySeriesService() *SeriesService {
	return NewSeriesService(DummyAffliateID, DummyAPIID)
}
