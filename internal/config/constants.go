package config

import (
	"time"

	"github.com/hogecode/JikkyoUtil/internal/models"
)

// JST is Japan Standard Time (UTC+9)
var JST = time.FixedZone("JST", 9*60*60)

// API endpoints
const (
	SyoboiTitleSearchURL = "http://cal.syoboi.jp/json"
	SyoboiProgLookupURL  = "http://cal.syoboi.jp/db"
	JikkyoBaseURL        = "https://jikkyo.tsukumijima.net/api/kakolog"
)

// NewChannelMapping creates and returns a complete channel mapping
func NewChannelMapping() models.ChannelMapping {
	return models.ChannelMapping{
		"NHK総合":    {ChID: 1, ChGID: 11, ChName: "NHK総合", JikkyoID: "jk1"},
		"NHK Eテレ":  {ChID: 2, ChGID: 11, ChName: "NHK Eテレ", JikkyoID: "jk2"},
		"フジテレビ":    {ChID: 3, ChGID: 1, ChName: "フジテレビ", JikkyoID: "jk8"},
		"日本テレビ":    {ChID: 4, ChGID: 1, ChName: "日テレ１", JikkyoID: "jk4"},
		"TBS":      {ChID: 5, ChGID: 1, ChName: "ＴＢＳ１", JikkyoID: "jk6"},
		"テレビ朝日":    {ChID: 6, ChGID: 1, ChName: "テレビ朝日", JikkyoID: "jk5"},
		"テレビ東京":    {ChID: 7, ChGID: 1, ChName: "テレビ東京１", JikkyoID: "jk7"},
		"tvk":      {ChID: 8, ChGID: 1, ChName: "tvk", JikkyoID: "jk11"},
		//"NHK-BS1":  {ChID: 9, ChGID: 9, ChName: "NHK-BS1", JikkyoID: "jk101"},
		//"NHK-BS2":  {ChID: 10, ChGID: 9, ChName: "NHK-BS2", JikkyoID: "jk101"},
		//"NHK-BShi": {ChID: 11, ChGID: 2, ChName: "NHK-BShi", JikkyoID: "jk103"},
		//"WOWOW":    {ChID: 12, ChGID: 9, ChName: "WOWOW", JikkyoID: "jk191"},
		//"チバテレビ":    {ChID: 13, ChGID: 1, ChName: "チバテレビ", JikkyoID: "jk12"},
		//"テレ玉":      {ChID: 14, ChGID: 1, ChName: "テレ玉", JikkyoID: "jk10"},
		//"BSテレ東":    {ChID: 15, ChGID: 2, ChName: "BSテレ東", JikkyoID: "jk171"},
		//"BS-TBS":   {ChID: 16, ChGID: 2, ChName: "BS-TBS", JikkyoID: "jk161"},
		//"BSフジ":     {ChID: 17, ChGID: 2, ChName: "BSフジ", JikkyoID: "jk181"},
		//"BS朝日":     {ChID: 18, ChGID: 2, ChName: "BS朝日", JikkyoID: "jk151"},
		"TOKYO MX": {ChID: 19, ChGID: 1, ChName: "ＴＯＫＹＯ　ＭＸ１", JikkyoID: "jk9"},
	}
}

// TimeFormat is the standard time format used for output
const TimeFormat = "2006-01-02 15:04:05"

// IsPopularChannel checks if a JikkyoID is for a popular/widely-supported channel (jk1-jk9)
// These channels have better comment availability and wider Jikkyo support
func IsPopularChannel(jikkyoID string) bool {
	// jk1 to jk9 are the main terrestrial/ground wave channels
	return jikkyoID == "jk1" || jikkyoID == "jk2" || jikkyoID == "jk4" ||
		jikkyoID == "jk5" || jikkyoID == "jk6" || jikkyoID == "jk7" ||
		jikkyoID == "jk8" || jikkyoID == "jk9"
}
