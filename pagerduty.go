package pagerduty

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Account struct {
	apiKey string
	url    string
}

type Incident struct {
	AssignedToUser        map[string]interface{} `json:"assigned_to_user"`
	CreatedOn             string                 `json:"created_on"`
	HtmlUrl               string                 `json:"html_url"`
	IncidentKey           string                 `json:"incident_key"`
	IncidentNumber        int                    `json:"incident_number"`
	LastStatusChangeOn    string                 `json:"last_status_change_on"`
	Service               map[string]interface{}
	Status                string
	TriggerDetailsHtmlUrl string                 `json:"trigger_details_html_url"`
	TriggerSummaryData    map[string]interface{} `json:"trigger_summary_data"`
}

type IncidentsResponse struct {
	Incidents []Incident
	Limit     int
	Offset    int
	Total     int
}

func SetupAccount(subdomain string, apiKey string) (account Account) {
	account = Account{apiKey: apiKey, url: fmt.Sprintf("https://%s.pagerduty.com", subdomain)}
	return
}

func (account *Account) Incidents(params map[string]string) (incidents []Incident, err error) {
	var (
		buf  []byte
		req  *http.Request
		resp *http.Response
	)

	endpoint := "api/v1/incidents"

	if len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Set(k, v)
		}
		endpoint = fmt.Sprintf("%s?%s", endpoint, values.Encode())
	}

	if req, err = account.getRequest(endpoint); err != nil {
		return
	}

	if resp, err = http.DefaultClient.Do(req); err != nil {
		return
	}

	if buf, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	defer resp.Body.Close()

	response := &IncidentsResponse{}

	if err = json.Unmarshal(buf, response); err != nil {
		return
	}

	incidents = response.Incidents
	return
}

func (account *Account) getRequest(endpoint string) (req *http.Request, err error) {
	if req, err = http.NewRequest("GET", fmt.Sprintf("%s/%s", account.url, endpoint), nil); err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Token token=%s", account.apiKey))

	return
}
