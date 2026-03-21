package api

import (
	"encoding/xml"
	"fmt"

	"github.com/hogecode/getabc/internal/config"
	"github.com/hogecode/getabc/internal/models"
)

// TitleSearch calls the Syoboi TitleSearch API
func (c *Client) TitleSearch(searchQuery string) (*models.TitleSearchResponse, error) {
	resp, err := c.R().
		SetQueryParams(map[string]string{
			"Req":    "TitleSearch",
			"Search": searchQuery,
			"Limit":  "15",
		}).
		SetResult(&models.TitleSearchResponse{}).
		Get(config.SyoboiTitleSearchURL)

	if err != nil {
		return nil, fmt.Errorf("failed to call TitleSearch API: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("TitleSearch API returned status %d", resp.StatusCode())
	}

	result, ok := resp.Result().(*models.TitleSearchResponse)
	if !ok {
		return nil, fmt.Errorf("failed to parse TitleSearch response")
	}

	return result, nil
}

// ProgLookup calls the Syoboi ProgLookup API
func (c *Client) ProgLookup(tid string, chID int, count int) (*models.ProgLookupResponse, error) {
	resp, err := c.R().
		SetQueryParams(map[string]string{
			"Command": "ProgLookup",
			"TID":     tid,
			"ChID":    fmt.Sprintf("%d", chID),
			"Count":   fmt.Sprintf("%d", count),
			"JOIN":    "SubTitles",
		}).
		Get(config.SyoboiProgLookupURL)

	if err != nil {
		return nil, fmt.Errorf("failed to call ProgLookup API: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("ProgLookup API returned status %d", resp.StatusCode())
	}

	// Parse XML response
	var result models.ProgLookupResponse
	if err := xml.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse ProgLookup response: %w", err)
	}

	return &result, nil
}
