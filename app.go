package steamapi

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func reqReadUnmarshal(url string, i interface{}, isXml bool) error {
	req, err := http.Get(url)
	if err != nil {
		return err
	}

	rawresp, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	if isXml {
		return xml.Unmarshal(rawresp, i)
	} else {
		return json.Unmarshal(rawresp, i)
	}
}

type AppList struct {
	XMLName xml.Name `xml:"applist"`
	AppList []struct {
		AppId int64  `xml:"appid"`
		Name  string `xml:"name"`
	} `xml:"apps>app"`
}

func GetAppList() (*AppList, error) {
	resp := &AppList{}
	return resp, reqReadUnmarshal("http://api.steampowered.com/ISteamApps/GetAppList/v0002/?format=xml&v=1", resp, true)
}

type AppDetailed struct {
	AboutTheGame string          `json:"about_the_game"`
	Platforms    map[string]bool `json:"platforms"`
	Categories   []struct {
		Description string `json:"description"`
		Id          int    `json:"id"`
	} `json:"categories"`
	Genres []struct {
		Description string `json:"description"`
		Id          string `json:"id"`
	} `json:"genres"`
	RequirementsPc      interface{} `json:"pc_requirements"`
	RequirementsMac     interface{} `json:"mac_requirements"`
	RequirementsLinux   interface{} `json:"linux_requirements"`
	RequiredAge         int         `json:"required_age"`
	IsFree              bool        `json:"is_free"`
	ControllerSupport   string      `json:"controller_support"`
	Background          string      `json:"background"`
	HeaderImage         string      `json:"header_image"`
	DetailedDescription string      `json:"detailed_description"`
	Publishers          []string    `json:"publishers"`
	Developers          []string    `json:"developers"`
	SupportedLanguages  string      `json:"supported_languages"`
	Type                string      `json:"type"`
	WebSite             string      `json:"website"`
	PriceOverview       struct {
		Currency        string `json:"currency"`
		DiscountPercent int    `json:"discount_percent"`
		Initial         int    `json:"initial"`
		Final           int    `json:"final"`
	} `json:"price_overview"`
}

type appDetailedWrapper map[string]struct {
	Data    *AppDetailed `json:"data"`
	Success bool         `json:"success"`
}

func GetAppDetailed(appid int64, cc, language string) (*AppDetailed, error) {
	var resp appDetailedWrapper
	url := fmt.Sprintf("http://store.steampowered.com/api/appdetails/?appids=%d&cc=%s&l=%s&v=1", appid, cc, language)
	if err := reqReadUnmarshal(url, &resp, false); err != nil {
		return nil, err
	}
	if r, ex := resp[strconv.FormatInt(appid, 10)]; ex {
		if !r.Success {
			return nil, errors.New("Steam tell unsuccess when receiving app details")
		} else {
			return r.Data, nil
		}
	} else {
		return nil, errors.New("Cannot find game with that id")
	}
}
