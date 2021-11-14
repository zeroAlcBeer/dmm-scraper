package api

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

type SeriesService struct {
	ApiID       string `mapstructure:"api_id"`
	AffiliateID string `mapstructure:"affiliate_id"`
	FloorID     string `mapstructure:"floor_id"`
	Initial     string `mapstructure:"initial"`
	Length      int64  `mapstructure:"hits"`
	Offset      int64  `mapstructure:"offset"`
}

type SeriesRawResponse struct {
	Request SeriesService  `mapstructure:"request"`
	Result  SeriesResponse `mapstructure:"result"`
}

type SeriesResponse struct {
	ResultCount   int64    `mapstructure:"result_count"`
	TotalCount    int64    `mapstructure:"total_count"`
	FirstPosition int64    `mapstructure:"first_position"`
	SiteName      string   `mapstructure:"site_name"`
	SiteCode      string   `mapstructure:"site_code"`
	ServiceName   string   `mapstructure:"service_name"`
	ServiceCode   string   `mapstructure:"service_code"`
	FloorID       string   `mapstructure:"floor_id"`
	FloorName     string   `mapstructure:"floor_name"`
	FloorCode     string   `mapstructure:"floor_code"`
	SeriesList    []Series `mapstructure:"series"`
}

type Series struct {
	SeriesID string `mapstructure:"series_id"`
	Name     string `mapstructure:"name"`
	Ruby     string `mapstructure:"ruby"`
	ListURL  string `mapstructure:"list_url"`
}

// NewSeriesService returns a new service for the given affiliate ID and API ID.
//
// NewSeriesServiceは渡したアフィリエイトIDとAPI IDを使用して新しい serviceを返します。
func NewSeriesService(affiliateID, apiID string) *SeriesService {
	return &SeriesService{
		ApiID:       apiID,
		AffiliateID: affiliateID,
		FloorID:     "",
		Initial:     "",
		Length:      DefaultAPILength,
		Offset:      DefaultAPIOffset,
	}
}

// Execute requests a url is created by BuildRequestURL.
// Use ExecuteWeak If you want get this response in interface{}.
//
// BuildRequestURLで生成したURLにリクエストします。
// もし interface{} でこのレスポンスを取得したい場合は ExecuteWeak を使用してください。
func (srv *SeriesService) Execute() (*SeriesResponse, error) {
	result, err := srv.ExecuteWeak()
	if err != nil {
		return nil, err
	}
	var raw SeriesRawResponse
	if err = mapstructure.WeakDecode(result, &raw); err != nil {
		return nil, err
	}
	return &raw.Result, nil
}

// ExecuteWeak requests a url is created by BuildRequestURL.
//
// BuildRequestURLで生成したURLにリクエストします。
func (srv *SeriesService) ExecuteWeak() (interface{}, error) {
	reqURL, err := srv.BuildRequestURL()
	if err != nil {
		return nil, err
	}

	return RequestJSON(reqURL)
}

// SetLength set the specified argument to SeriesService.Length
//
// SetLengthはLengthパラメータを設定します。
func (srv *SeriesService) SetLength(length int64) *SeriesService {
	srv.Length = length
	return srv
}

// SetHits set the specified argument to SeriesService.Length
//  SetHits is the alias for SetLength
//
// SetHitsはLengthパラメータを設定します。
func (srv *SeriesService) SetHits(length int64) *SeriesService {
	srv.SetLength(length)
	return srv
}

// SetOffset set the specified argument to SeriesService.Offset
//
// SetOffsetはOffsetパラメータを設定します。
func (srv *SeriesService) SetOffset(offset int64) *SeriesService {
	srv.Offset = offset
	return srv
}

// SetInitial sets the specified argument to SeriesService.Initial.
// This argment is series' initial and you can use only hiragana.
//  e.g. srv.SetInitial("あ") -> ARIA(ありあ)
//
// SetInitialはInitalパラメータに検索したいシリーズの頭文字をひらがなで設定します。
func (srv *SeriesService) SetInitial(initial string) *SeriesService {
	srv.Initial = TrimString(initial)
	return srv
}

// SetFloorID sets the specified argument to SeriesService.FloorID.
// You can retrieve Floor IDs from floor API.
//
// SetFloorIDはFloorIDパラメータを設定します。
// フロアIDはフロアAPIから取得できます。
func (srv *SeriesService) SetFloorID(floorID string) *SeriesService {
	srv.FloorID = TrimString(floorID)
	return srv
}

// ValidateLength validates SeriesService.Length within the range (1 <= value <= DefaultMaxLength).
// Refer to ValidateRange for more information about the range to validate.
//
// ValidateLengthはSeriesService.Lengthが範囲内(1 <= value <= DefaultMaxLength)にあるか検証します。
// 検証範囲について更に詳しく知りたい方はValidateRangeを参照してください。
func (srv *SeriesService) ValidateLength() bool {
	return ValidateRange(srv.Length, 1, DefaultMaxLength)
}

// ValidateOffset validates SeriesService.Offset within the range (1 <= value).
//
// ValidateOffsetはSeriesService.Offsetが範囲内(1 <= value)にあるか検証します。
func (srv *SeriesService) ValidateOffset() bool {
	return srv.Offset >= 1
}

// BuildRequestURL creates url to request series API.
//
// BuildRequestURLはシリーズ検索APIにリクエストするためのURLを作成します。
func (srv *SeriesService) BuildRequestURL() (string, error) {
	if srv.ApiID == "" {
		return "", fmt.Errorf("set invalid ApiID parameter")
	}
	if !ValidateAffiliateID(srv.AffiliateID) {
		return "", fmt.Errorf("set invalid AffiliateID parameter")
	}
	if srv.FloorID == "" {
		return "", fmt.Errorf("set invalid FloorID parameter")
	}

	queries := url.Values{}
	queries.Set("api_id", srv.ApiID)
	queries.Set("affiliate_id", srv.AffiliateID)
	queries.Set("floor_id", srv.FloorID)

	if srv.Length != 0 {
		if !srv.ValidateLength() {
			return "", fmt.Errorf("length out of range: %d", srv.Length)
		}
		queries.Set("hits", strconv.FormatInt(srv.Length, 10))
	}

	if srv.Offset != 0 {
		if !srv.ValidateOffset() {
			return "", fmt.Errorf("offset out of range: %d", srv.Offset)
		}
		queries.Set("offset", strconv.FormatInt(srv.Offset, 10))
	}

	if srv.Initial != "" {
		queries.Set("initial", srv.Initial)
	}

	u, err := buildAPIEndpoint("SeriesSearch")
	if err != nil {
		return "", err
	}
	u.RawQuery = queries.Encode()

	return u.String(), nil
}
