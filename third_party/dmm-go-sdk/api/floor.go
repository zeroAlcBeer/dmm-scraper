package api

import (
	"fmt"
	"net/url"

	"github.com/mitchellh/mapstructure"
)

type FloorService struct {
	ApiID       string `mapstructure:"api_id"`
	AffiliateID string `mapstructure:"affiliate_id"`
}

type FloorRawResponse struct {
	Request FloorService  `mapstructure:"request"`
	Result  FloorResponse `mapstructure:"result"`
}

type FloorResponse struct {
	Site []Site
}

type Site struct {
	Name     string       `mapstructure:"name"`
	Code     string       `mapstructure:"code"`
	Services []DMMService `mapstructure:"service"`
}

type DMMService struct {
	Name   string     `mapstructure:"name"`
	Code   string     `mapstructure:"code"`
	Floors []DMMFloor `mapstructure:"floor"`
}

type DMMFloor struct {
	ID   int64  `mapstructure:"id"`
	Name string `mapstructure:"name"`
	Code string `mapstructure:"code"`
}

// NewFloorService returns a new service for the given affiliate ID and API ID.
//
// NewFloorServiceは渡したアフィリエイトIDとAPI IDを使用して新しい serviceを返します。
func NewFloorService(affiliateID, apiID string) *FloorService {
	return &FloorService{
		ApiID:       apiID,
		AffiliateID: affiliateID,
	}
}

// Execute requests a url is created by BuildRequestURL.
// Use ExecuteWeak If you want get this response in interface{}.
//
// BuildRequestURLで生成したURLにリクエストします。
// もし interface{} でこのレスポンスを取得したい場合は ExecuteWeak を使用してください。
func (srv *FloorService) Execute() (*FloorResponse, error) {
	result, err := srv.ExecuteWeak()
	if err != nil {
		return nil, err
	}
	var raw FloorRawResponse
	if err = mapstructure.WeakDecode(result, &raw); err != nil {
		return nil, err
	}
	return &raw.Result, nil
}

// ExecuteWeak requests a url is created by BuildRequestURL.
//
// BuildRequestURLで生成したURLにリクエストします。
func (srv *FloorService) ExecuteWeak() (interface{}, error) {
	reqURL, err := srv.BuildRequestURL()
	if err != nil {
		return nil, err
	}

	return RequestJSON(reqURL)
}

// BuildRequestURL creates url to request floor API.
//
// BuildRequestURLはフロアAPIにリクエストするためのURLを作成します。
func (srv *FloorService) BuildRequestURL() (string, error) {
	if srv.ApiID == "" {
		return "", fmt.Errorf("set invalid ApiID parameter")
	}
	if !ValidateAffiliateID(srv.AffiliateID) {
		return "", fmt.Errorf("set invalid AffiliateID parameter")
	}

	queries := url.Values{}
	queries.Set("api_id", srv.ApiID)
	queries.Set("affiliate_id", srv.AffiliateID)

	u, err := buildAPIEndpoint("FloorList")
	if err != nil {
		return "", err
	}
	u.RawQuery = queries.Encode()

	return u.String(), nil
}
