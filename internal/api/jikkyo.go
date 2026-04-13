package api

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/hogecode/JikkyoUtil/internal/config"
	"github.com/hogecode/JikkyoUtil/internal/models"
)

// GetJikkyoComments fetches comments from Jikkyo API
func (c *Client) GetJikkyoComments(jikkyoID string, startTime, endTime int64) (*models.JikkyoResponse, error) {
	url := fmt.Sprintf("%s/%s", config.JikkyoBaseURL, jikkyoID)

	resp, err := c.R().
		SetQueryParams(map[string]string{
			"starttime": fmt.Sprintf("%d", startTime),
			"endtime":   fmt.Sprintf("%d", endTime),
			"format":    "json",
		}).
		SetResult(&models.JikkyoResponse{}).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to call Jikkyo API: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("Jikkyo API returned status %d", resp.StatusCode())
	}

	result, ok := resp.Result().(*models.JikkyoResponse)
	if !ok {
		return nil, fmt.Errorf("failed to parse Jikkyo response")
	}

	return result, nil
}

// GetJikkyoCommentsXML fetches comments from Jikkyo API in XML format
func (c *Client) GetJikkyoCommentsXML(jikkyoID string, startTime, endTime int64) (*models.Packet, error) {
	url := fmt.Sprintf("%s/%s", config.JikkyoBaseURL, jikkyoID)

	resp, err := c.R().
		SetQueryParams(map[string]string{
			"starttime": fmt.Sprintf("%d", startTime),
			"endtime":   fmt.Sprintf("%d", endTime),
			"format":    "xml",
		}).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to call Jikkyo API: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("Jikkyo API returned status %d", resp.StatusCode())
	}

	var result models.Packet
	err = xml.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Jikkyo XML response: %w", err)
	}

	// Normalize vpos values: set the first comment's vpos to 0 and adjust all others accordingly
	if len(result.Chats) > 0 {
		firstVpos, err := strconv.ParseInt(result.Chats[0].Vpos, 10, 64)
		if err == nil && firstVpos > 0 {
			for i := range result.Chats {
				currentVpos, err := strconv.ParseInt(result.Chats[i].Vpos, 10, 64)
				if err == nil {
					normalizedVpos := currentVpos - firstVpos
					result.Chats[i].Vpos = fmt.Sprintf("%d", normalizedVpos)
				}
			}
		}
	}

	return &result, nil
}
