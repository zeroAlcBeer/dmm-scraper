package api

import (
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

func TestNewMakerService(t *testing.T) {
	affiliateID := DummyAffliateID
	apiID := DummyAPIID

	srv := NewMakerService(affiliateID, apiID)
	if srv.AffiliateID != affiliateID {
		t.Fatalf("MakerService.AffiliateID is expected to equal the input value(affiliateID)")
	}

	if srv.ApiID != apiID {
		t.Fatalf("MakerService.ApiID is expected to equal the input value(apiID)")
	}
}

func TestSetLengthInMakerService(t *testing.T) {
	srv := dummyMakerService()
	var length int64 = 10
	srv.SetLength(length)

	if srv.Length != length {
		t.Fatalf("MakerService.Length is expected to equal the input value(length)")
	}
}

func TestSetHitsInMakerService(t *testing.T) {
	srv := dummyMakerService()
	var hits int64 = 10
	srv.SetHits(hits)

	if srv.Length != hits {
		t.Fatalf("MakerService.Length is expected to equal the input value(hits)")
	}
}

func TestSetOffsetInMakerService(t *testing.T) {
	srv := dummyMakerService()
	var offset int64 = 10
	srv.SetOffset(offset)

	if srv.Offset != offset {
		t.Fatalf("MakerService.Offset is expected to equal the input value(offset)")
	}
}

func TestValidateLengthInMakerService(t *testing.T) {
	srv := dummyMakerService()

	var target int64

	target = 1
	srv.SetLength(target)
	if srv.ValidateLength() == false {
		t.Fatalf("MakerService.ValidateLength is expected TRUE.")
	}

	target = DefaultAPILength
	srv.SetLength(target)
	if srv.ValidateLength() == false {
		t.Fatalf("MakerService.ValidateLength is expected TRUE.")
	}

	target = DefaultMaxLength
	srv.SetLength(target)
	if srv.ValidateLength() == false {
		t.Fatalf("MakerService.ValidateLength is expected TRUE.")
	}

	target = DefaultMaxLength + 1
	srv.SetLength(target)
	if srv.ValidateLength() == true {
		t.Fatalf("MakerService.ValidateLength is expected FALSE.")
	}

	target = 0
	srv.SetLength(target)
	if srv.ValidateLength() == true {
		t.Fatalf("MakerService.ValidateLength is expected FALSE.")
	}

	target = -1
	srv.SetLength(target)
	if srv.ValidateLength() == true {
		t.Fatalf("MakerService.ValidateLength is expected FALSE")
	}
}

func TestValidateOffsetInMakerService(t *testing.T) {
	srv := dummyMakerService()

	var target int64

	target = 1
	srv.SetOffset(target)
	if srv.ValidateOffset() == false {
		t.Fatalf("MakerService.ValidateOffset is expected TRUE")
	}

	target = 100
	srv.SetOffset(target)
	if srv.ValidateOffset() == false {
		t.Fatalf("MakerService.ValidateOffset is expected TRUE")
	}

	target = 0
	srv.SetOffset(target)
	if srv.ValidateOffset() == true {
		t.Fatalf("MakerService.ValidateOffset is expected FALSE")
	}

	target = -1
	srv.SetOffset(target)
	if srv.ValidateOffset() == true {
		t.Fatalf("MakerService.ValidateOffset is expected FALSE")
	}
}

func TestBuildRequestURLInMakerService(t *testing.T) {
	var srv *MakerService
	var u string
	var err error
	var expected string

	srv = dummyMakerService()
	srv.SetFloorID("40")
	u, err = srv.BuildRequestURL()
	expected = APIBaseURL + "/MakerSearch?affiliate_id=" + DummyAffliateID + "&api_id=" + DummyAPIID + "&floor_id=40" + "&hits=" + strconv.FormatInt(DefaultAPILength, 10) + "&offset=" + strconv.FormatInt(DefaultAPIOffset, 10)
	if u != expected {
		t.Fatalf("MakerService.BuildRequestURL is expected to equal the expected value.\nexpected:%s\nactual:  %s", expected, u)
	}
	if err != nil {
		t.Fatalf("MakerService.BuildRequestURL is not expected to have error")
	}

	srv = dummyMakerService()
	srv.SetLength(0)
	srv.SetOffset(0)
	srv.SetFloorID("40")
	u, err = srv.BuildRequestURL()
	expected = APIBaseURL + "/MakerSearch?affiliate_id=" + DummyAffliateID + "&api_id=" + DummyAPIID + "&floor_id=40"
	expectedBase := expected
	if u != expected {
		t.Fatalf("MakerService.BuildRequestURL is expected to equal the expected value.\nexpected:%s\nactual:  %s", expected, u)
	}
	if err != nil {
		t.Fatalf("MakerService.BuildRequestURL is not expected to have error")
	}

	srv.SetLength(-1)
	u, err = srv.BuildRequestURL()
	if u != "" {
		t.Fatalf("MakerService.BuildRequestURL is expected empty if error occurs.")
	}
	if err == nil {
		t.Fatalf("MakerService.BuildRequestURL is expected to return error.")
	}
	srv.SetLength(0)

	srv.SetOffset(-1)
	u, err = srv.BuildRequestURL()
	if u != "" {
		t.Fatalf("MakerService.BuildRequestURL is expected empty if error occurs.")
	}
	if err == nil {
		t.Fatalf("MakerService.BuildRequestURL is expected to return error.")
	}
	srv.SetOffset(0)

	srv.SetFloorID("")
	u, err = srv.BuildRequestURL()
	if u != "" {
		t.Fatalf("MakerService.BuildRequestURL is expected empty if error occurs.")
	}
	if err == nil {
		t.Fatalf("MakerService.BuildRequestURL is expected to return error.")
	}
	srv.SetFloorID("40")

	srv.SetInitial("あ")
	expected = expectedBase + "&initial=" + url.QueryEscape("あ")
	u, err = srv.BuildRequestURL()
	if u != expected {
		t.Fatalf("MakerService.BuildRequestURL is expected to equal the expected value.\nexpected:%s\nactual:  %s", expected, u)
	}
	if err != nil {
		t.Fatalf("MakerService.BuildRequestURL is not expected to have error")
	}
	srv.SetInitial("")
}

func TestBuildRequestURLWithoutApiIDInMakerService(t *testing.T) {
	srv := dummyMakerService()
	srv.ApiID = ""
	u, err := srv.BuildRequestURL()
	if u != "" {
		t.Fatalf("MakerService.BuildRequestURL is expected empty if API ID is not set.")
	}
	if err == nil {
		t.Fatalf("MakerService.BuildRequestURL is expected to return error.")
	}
}

func TestBuildRequestURLWithWrongAffiliateIDInMakerService(t *testing.T) {
	srv := dummyMakerService()
	srv.AffiliateID = "fizzbizz-100"
	u, err := srv.BuildRequestURL()
	if u != "" {
		t.Fatalf("MakerService.BuildRequestURL is expected empty if wrong Affiliate ID is set.")
	}
	if err == nil {
		t.Fatalf("MakerService.BuildRequestURL is expected to return error.")
	}
}

func TestExcuteRequestMakerAPIToServer(t *testing.T) {
	if !RequestAvailable {
		t.Skip("Not set valid credentials")
	}

	srv := NewMakerService(TestAffiliateID, TestAPIID)
	srv.SetFloorID("40")
	srv.SetInitial("あ")
	srv.SetLength(100)
	srv.SetOffset(1)

	rst, err := srv.Execute()
	if err != nil {
		t.Skip("Maybe, The network is down.")
	}

	if reflect.TypeOf(rst).String() != "*api.MakerResponse" {
		t.Fatalf("MakerService.Execute is expected to return *api.MakerResponse")
	}
}

func dummyMakerService() *MakerService {
	return NewMakerService(DummyAffliateID, DummyAPIID)
}
