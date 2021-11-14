package api

import (
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

func TestNewAuthorService(t *testing.T) {
	affiliateID := DummyAffliateID
	apiID := DummyAPIID

	srv := NewAuthorService(affiliateID, apiID)
	if srv.AffiliateID != affiliateID {
		t.Fatalf("AuthorService.AffiliateID is expected to equal the input value(affiliateID)")
	}

	if srv.ApiID != apiID {
		t.Fatalf("AuthorService.ApiID is expected to equal the input value(apiID)")
	}
}

func TestSetLengthInAuthorService(t *testing.T) {
	srv := dummyAuthorService()
	var length int64 = 10
	srv.SetLength(length)

	if srv.Length != length {
		t.Fatalf("AuthorService.Length is expected to equal the input value(length)")
	}
}

func TestSetHitsInAuthorService(t *testing.T) {
	srv := dummyAuthorService()
	var hits int64 = 10
	srv.SetHits(hits)

	if srv.Length != hits {
		t.Fatalf("AuthorService.Length is expected to equal the input value(hits)")
	}
}

func TestSetOffsetInAuthorService(t *testing.T) {
	srv := dummyAuthorService()
	var offset int64 = 10
	srv.SetOffset(offset)

	if srv.Offset != offset {
		t.Fatalf("AuthorService.Offset is expected to equal the input value(offset)")
	}
}

func TestValidateLengthInAuthorService(t *testing.T) {
	srv := dummyAuthorService()

	var target int64

	target = 1
	srv.SetLength(target)
	if srv.ValidateLength() == false {
		t.Fatalf("AuthorService.ValidateLength is expected TRUE.")
	}

	target = DefaultAPILength
	srv.SetLength(target)
	if srv.ValidateLength() == false {
		t.Fatalf("AuthorService.ValidateLength is expected TRUE.")
	}

	target = DefaultMaxLength
	srv.SetLength(target)
	if srv.ValidateLength() == false {
		t.Fatalf("AuthorService.ValidateLength is expected TRUE.")
	}

	target = DefaultMaxLength + 1
	srv.SetLength(target)
	if srv.ValidateLength() == true {
		t.Fatalf("AuthorService.ValidateLength is expected FALSE.")
	}

	target = 0
	srv.SetLength(target)
	if srv.ValidateLength() == true {
		t.Fatalf("AuthorService.ValidateLength is expected FALSE.")
	}

	target = -1
	srv.SetLength(target)
	if srv.ValidateLength() == true {
		t.Fatalf("AuthorService.ValidateLength is expected FALSE.")
	}
}

func TestValidateOffsetInAuthorService(t *testing.T) {
	srv := dummyAuthorService()

	var target int64

	target = 1
	srv.SetOffset(target)
	if srv.ValidateOffset() == false {
		t.Fatalf("AuthorService.ValidateOffset is expected TRUE.")
	}

	target = 100
	srv.SetOffset(target)
	if srv.ValidateOffset() == false {
		t.Fatalf("AuthorService.ValidateOffset is expected TRUE.")
	}

	target = 0
	srv.SetOffset(target)
	if srv.ValidateOffset() == true {
		t.Fatalf("AuthorService.ValidateOffset is expected FALSE.")
	}

	target = -1
	srv.SetOffset(target)
	if srv.ValidateOffset() == true {
		t.Fatalf("AuthorService.ValidateOffset is expected FALSE.")
	}
}

func TestBuildRequestURLInAuthorService(t *testing.T) {
	var srv *AuthorService
	var u string
	var err error
	var expected string

	srv = dummyAuthorService()
	srv.SetFloorID("40")
	u, err = srv.BuildRequestURL()
	expected = APIBaseURL + "/AuthorSearch?affiliate_id=" + DummyAffliateID + "&api_id=" + DummyAPIID + "&floor_id=40" + "&hits=" + strconv.FormatInt(DefaultAPILength, 10) + "&offset=" + strconv.FormatInt(DefaultAPIOffset, 10)
	if u != expected {
		t.Fatalf("AuthorService.BuildRequestURL is expected to equal the expected value.\nexpected:%s\nactual:  %s", expected, u)
	}
	if err != nil {
		t.Fatalf("AuthorService.BuildRequestURL is not expected to have error")
	}

	srv = dummyAuthorService()
	srv.SetLength(0)
	srv.SetOffset(0)
	srv.SetFloorID("40")
	u, err = srv.BuildRequestURL()
	expected = APIBaseURL + "/AuthorSearch?affiliate_id=" + DummyAffliateID + "&api_id=" + DummyAPIID + "&floor_id=40"
	expectedBase := expected
	if u != expected {
		t.Fatalf("AuthorService.BuildRequestURL is expected to equal the expected value.\nexpected:%s\nactual:  %s", expected, u)
	}
	if err != nil {
		t.Fatalf("AuthorService.BuildRequestURL is not expected to have error")
	}

	srv.SetLength(-1)
	u, err = srv.BuildRequestURL()
	if u != "" {
		t.Fatalf("AuthorService.BuildRequestURL is expected empty if error occurs.")
	}
	if err == nil {
		t.Fatalf("AuthorService.BuildRequestURL is expected to return error.")
	}
	srv.SetLength(0)

	srv.SetOffset(-1)
	u, err = srv.BuildRequestURL()
	if u != "" {
		t.Fatalf("AuthorService.BuildRequestURL is expected empty if error occurs.")
	}
	if err == nil {
		t.Fatalf("AuthorService.BuildRequestURL is expected to return error.")
	}
	srv.SetOffset(0)

	srv.SetFloorID("")
	u, err = srv.BuildRequestURL()
	if u != "" {
		t.Fatalf("AuthorService.BuildRequestURL is expected empty if error occurs.")
	}
	if err == nil {
		t.Fatalf("AuthorService.BuildRequestURL is expected to return error.")
	}
	srv.SetFloorID("40")

	srv.SetInitial("あ")
	expected = expectedBase + "&initial=" + url.QueryEscape("あ")
	u, err = srv.BuildRequestURL()
	if u != expected {
		t.Fatalf("AuthorService.BuildRequestURL is expected to equal the expected value.\nexpected:%s\nactual:  %s", expected, u)
	}
	if err != nil {
		t.Fatalf("AuthorService.BuildRequestURL is not expected to have error")
	}
	srv.SetInitial("")
}

func TestBuildRequestURLWithoutApiIDInAuthorService(t *testing.T) {
	srv := dummyAuthorService()
	srv.ApiID = ""
	u, err := srv.BuildRequestURL()
	if u != "" {
		t.Fatalf("AuthorService.BuildRequestURL is expected empty if API ID is not set.")
	}
	if err == nil {
		t.Fatalf("AuthorService.BuildRequestURL is expected to return error.")
	}
}

func TestBuildRequestURLWithWrongAffiliateIDInAuthorService(t *testing.T) {
	srv := dummyAuthorService()
	srv.AffiliateID = "fizzbizz-100"
	u, err := srv.BuildRequestURL()
	if u != "" {
		t.Fatalf("AuthorService.BuildRequestURL is expected empty if wrong Affiliate ID is set.")
	}
	if err == nil {
		t.Fatalf("AuthorService.BuildRequestURL is expected to return error.")
	}
}

func TestExcuteRequestAuthorAPIToServer(t *testing.T) {
	if !RequestAvailable {
		t.Skip("Not set valid credentials")
	}

	srv := NewAuthorService(TestAffiliateID, TestAPIID)
	srv.SetFloorID("40")
	srv.SetInitial("あ")
	srv.SetLength(100)
	srv.SetOffset(1)

	rst, err := srv.Execute()
	if err != nil {
		t.Skip("Maybe, The network is down.")
	}

	if reflect.TypeOf(rst).String() != "*api.AuthorResponse" {
		t.Fatalf("AuthorService.Execute is expected to return *api.AuthorResponse")
	}
}

func dummyAuthorService() *AuthorService {
	return NewAuthorService(DummyAffliateID, DummyAPIID)
}
