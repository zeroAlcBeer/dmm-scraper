package api

import (
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

func TestNewGenreService(t *testing.T) {
	affiliateID := DummyAffliateID
	apiID := DummyAPIID

	srv := NewGenreService(affiliateID, apiID)
	if srv.AffiliateID != affiliateID {
		t.Fatalf("GenreService.AffiliateID is expected to equal the input value(affiliateID)")
	}

	if srv.ApiID != apiID {
		t.Fatalf("GenreService.ApiID is expected to equal the input value(apiID)")
	}
}

func TestSetLengthInGenreService(t *testing.T) {
	srv := dummyGenreService()
	var length int64 = 10
	srv.SetLength(length)

	if srv.Length != length {
		t.Fatalf("GenreService.Length is expected to equal the input value(length)")
	}
}

func TestSetHitsInGenreService(t *testing.T) {
	srv := dummyGenreService()
	var hits int64 = 10
	srv.SetHits(hits)

	if srv.Length != hits {
		t.Fatalf("GenreService.Length is expected to equal the input value(hits)")
	}
}

func TestSetOffsetInGenreService(t *testing.T) {
	srv := dummyGenreService()
	var offset int64 = 10
	srv.SetOffset(offset)

	if srv.Offset != offset {
		t.Fatalf("GenreService.Offset is expected to equal the input value(offset)")
	}
}

func TestValidateLengthInGenreService(t *testing.T) {
	srv := dummyGenreService()

	var target int64

	target = 1
	srv.SetLength(target)
	if srv.ValidateLength() == false {
		t.Fatalf("GenreService.ValidateLength is expected TRUE.")
	}

	target = DefaultAPILength
	srv.SetLength(target)
	if srv.ValidateLength() == false {
		t.Fatalf("GenreService.ValidateLength is expected TRUE.")
	}

	target = DefaultMaxLength
	srv.SetLength(target)
	if srv.ValidateLength() == false {
		t.Fatalf("GenreService.ValidateLength is expected TRUE.")
	}

	target = DefaultMaxLength + 1
	srv.SetLength(target)
	if srv.ValidateLength() == true {
		t.Fatalf("GenreService.ValidateLength is expected FALSE.")
	}

	target = 0
	srv.SetLength(target)
	if srv.ValidateLength() == true {
		t.Fatalf("GenreService.ValidateLength is expected FALSE.")
	}

	target = -1
	srv.SetLength(target)
	if srv.ValidateLength() == true {
		t.Fatalf("GenreService.ValidateLength is expected FALSE.")
	}
}

func TestValidateOffsetInGenreService(t *testing.T) {
	srv := dummyGenreService()

	var target int64

	target = 1
	srv.SetOffset(target)
	if srv.ValidateOffset() == false {
		t.Fatalf("GenreService.ValidateOffset is expected TRUE.")
	}

	target = 100
	srv.SetOffset(target)
	if srv.ValidateOffset() == false {
		t.Fatalf("GenreService.ValidateOffset is expected TRUE.")
	}

	target = 0
	srv.SetOffset(target)
	if srv.ValidateOffset() == true {
		t.Fatalf("GenreService.ValidateOffset is expected FALSE.")
	}

	target = -1
	srv.SetOffset(target)
	if srv.ValidateOffset() == true {
		t.Fatalf("GenreService.ValidateOffset is expected FALSE.")
	}
}

func TestBuildRequestURLInGenreService(t *testing.T) {
	var srv *GenreService
	var u string
	var err error
	var expected string

	srv = dummyGenreService()
	srv.SetFloorID("40")
	u, err = srv.BuildRequestURL()
	expected = APIBaseURL + "/GenreSearch?affiliate_id=" + DummyAffliateID + "&api_id=" + DummyAPIID + "&floor_id=40" + "&hits=" + strconv.FormatInt(DefaultAPILength, 10) + "&offset=" + strconv.FormatInt(DefaultAPIOffset, 10)
	if u != expected {
		t.Fatalf("GenreService.BuildRequestURL is expected to equal the expected value.\nexpected:%s\nactual:  %s", expected, u)
	}
	if err != nil {
		t.Fatalf("GenreService.BuildRequestURL is not expected to have error")
	}

	srv = dummyGenreService()
	srv.SetLength(0)
	srv.SetOffset(0)
	srv.SetFloorID("40")
	u, err = srv.BuildRequestURL()
	expected = APIBaseURL + "/GenreSearch?affiliate_id=" + DummyAffliateID + "&api_id=" + DummyAPIID + "&floor_id=40"
	expectedBase := expected
	if u != expected {
		t.Fatalf("GenreService.BuildRequestURL is expected to equal the expected value.\nexpected:%s\nactual:  %s", expected, u)
	}
	if err != nil {
		t.Fatalf("GenreService.BuildRequestURL is not expected to have error")
	}

	srv.SetLength(-1)
	u, err = srv.BuildRequestURL()
	if u != "" {
		t.Fatalf("GenreService.BuildRequestURL is expected empty if error occurs.")
	}
	if err == nil {
		t.Fatalf("GenreService.BuildRequestURL is expected to return error.")
	}
	srv.SetLength(0)

	srv.SetOffset(-1)
	u, err = srv.BuildRequestURL()
	if u != "" {
		t.Fatalf("GenreService.BuildRequestURL is expected empty if error occurs.")
	}
	if err == nil {
		t.Fatalf("GenreService.BuildRequestURL is expected to return error.")
	}
	srv.SetOffset(0)

	srv.SetFloorID("")
	u, err = srv.BuildRequestURL()
	if u != "" {
		t.Fatalf("GenreService.BuildRequestURL is expected empty if error occurs.")
	}
	if err == nil {
		t.Fatalf("GenreService.BuildRequestURL is expected to return error.")
	}
	srv.SetFloorID("40")

	srv.SetInitial("あ")
	expected = expectedBase + "&initial=" + url.QueryEscape("あ")
	u, err = srv.BuildRequestURL()
	if u != expected {
		t.Fatalf("GenreService.BuildRequestURL is expected to equal the expected value.\nexpected:%s\nactual:  %s", expected, u)
	}
	if err != nil {
		t.Fatalf("GenreService.BuildRequestURL is not expected to have error")
	}
	srv.SetInitial("")
}

func TestBuildRequestURLWithoutApiIDInGenreService(t *testing.T) {
	srv := dummyGenreService()
	srv.ApiID = ""
	u, err := srv.BuildRequestURL()
	if u != "" {
		t.Fatalf("GenreService.BuildRequestURL is expected empty if API ID is not set.")
	}
	if err == nil {
		t.Fatalf("GenreService.BuildRequestURL is expected to return error.")
	}
}

func TestBuildRequestURLWithWrongAffiliateIDInGenreService(t *testing.T) {
	srv := dummyGenreService()
	srv.AffiliateID = "fizzbizz-100"
	u, err := srv.BuildRequestURL()
	if u != "" {
		t.Fatalf("GenreService.BuildRequestURL is expected empty if wrong Affiliate ID is set.")
	}
	if err == nil {
		t.Fatalf("GenreService.BuildRequestURL is expected to return error.")
	}
}

func TestExcuteRequestAPIServer(t *testing.T) {
	if !RequestAvailable {
		t.Skip("Not set valid credentials")
	}
}

func TestExcuteRequestGenreAPIToServer(t *testing.T) {
	if !RequestAvailable {
		t.Skip("Not set valid credentials")
	}

	srv := NewGenreService(TestAffiliateID, TestAPIID)
	srv.SetFloorID("40")
	srv.SetInitial("あ")
	srv.SetLength(100)
	srv.SetOffset(1)

	rst, err := srv.Execute()
	if err != nil {
		t.Skip("Maybe, The network is down.")
	}

	if reflect.TypeOf(rst).String() != "*api.GenreResponse" {
		t.Fatalf("GenreService.Execute is expected to return *api.GenreResponse")
	}
}

func dummyGenreService() *GenreService {
	return NewGenreService(DummyAffliateID, DummyAPIID)
}
