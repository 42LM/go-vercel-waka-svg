// Package service implements the actual service functions that deliver svg files.
package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

// WakaTimeInput type stores the information of interest
// from the waka time api last_7_days endpoint.
//
// This is used in the svg template to populate it with dynamic content.
type WakaTimeInput struct {
	Language string
	Time     string
	Bar      string
	Percent  string
	Y        int
}

// WakaResponse type references the response of the wakatime api.
type WakaResponse struct {
	Data struct {
		Languages []struct {
			Name         string  `json:"name"`
			TotalSeconds float64 `json:"total_seconds"`
			Percent      float64 `json:"percent"`
			Digital      string  `json:"digital"`
			Decimal      string  `json:"decimal"`
			Text         string  `json:"text"`
			Hours        int     `json:"hours"`
			Minutes      int     `json:"minutes"`
		} `json:"languages"`
	} `json:"data"`
}

// Wakatime returns the svg template wakatime and displays the last 7 days of coding activity.
// Templates can be found in `../svgtemplate/templates/...`
func (s *service) Wakatime(ctx context.Context) error {
	req, err := http.NewRequest(http.MethodGet, "https://wakatime.com/api/v1/users/current/stats/last_7_days", nil)
	if err != nil {
		return err
	}
	apiKey := os.Getenv("WAKA_API_KEY")
	base64APIKey := base64.StdEncoding.EncodeToString([]byte(apiKey))

	req.Header.Add("authorization", "Basic "+base64APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	wakaResp := WakaResponse{}
	err = json.Unmarshal(body, &wakaResp)
	if err != nil {
		return err
	}

	yCount := 90
	yCountAddidtion := 40
	// loop through all languages and check if it contains `Other`
	// set a different Ycount when that is the case
	var specialI *int
	for i, l := range wakaResp.Data.Languages {
		if l.Name == "Other" {
			specialI = &i
			yCountAddidtion = 55
			break
		}
	}

	// get the top 4 programming languages
	wakaT := make([]WakaTimeInput, 5)
	for i, l := range wakaResp.Data.Languages {
		if i == 4 {
			break
		}
		if l.Name == "Other" {
			continue
		}

		bar := calcBar(l.Percent)

		if i > 0 {
			yCount += yCountAddidtion
		}

		wakaT[i] = WakaTimeInput{
			Language: l.Name,
			Time:     l.Text,
			Bar:      bar,
			Percent:  fmt.Sprintf("%05.2f %%", l.Percent),
			Y:        yCount,
		}
	}

	// calculate all the rest
	other := WakaTimeInput{}
	var tmpPercent float64
	var tmpHours int
	var tmpMinutes int
	for i := 4; i < len(wakaResp.Data.Languages); i++ {
		l := wakaResp.Data.Languages[i]

		tmpPercent += l.Percent
		tmpHours += l.Hours
		tmpMinutes += l.Minutes

		if len(wakaResp.Data.Languages)-1 == i {
			bar := calcBar(tmpPercent)
			other.Bar = bar
		}
	}

	other.Language = "Other"

	// produce time text
	otherTime := ""
	hours, remainingMinutes := convertMinutesToHours(tmpMinutes)
	tmpHours += hours
	tmpMinutes = remainingMinutes

	if tmpHours > 0 && tmpMinutes != 0 {
		if tmpHours == 1 {
			otherTime = fmt.Sprintf("%d hr %d mins", tmpHours, tmpMinutes)
		} else {
			otherTime = fmt.Sprintf("%d hrs %d mins", tmpHours, tmpMinutes)
		}
	} else if tmpHours == 0 && tmpMinutes == 0 {
		otherTime = "0 secs"
	} else {
		otherTime = fmt.Sprintf("%d mins", tmpMinutes)
	}

	other.Time = otherTime
	other.Percent = fmt.Sprintf("%05.2f %%", tmpPercent)
	other.Y = yCount + yCountAddidtion

	if specialI != nil {
		tmpPercent += wakaResp.Data.Languages[*specialI].Percent
	}

	other.Percent = fmt.Sprintf("%05.2f %%", tmpPercent)
	wakaT[4] = other

	return s.templates.ExecuteTemplate(s.responseWriter, "wakatime.gosvg", wakaT)
}

// calcBar calculates the progress bar for each language.
func calcBar(percent float64) string {
	fullBars := (percent / 100) * 24
	restBars := 24 - fullBars
	truncated := math.Trunc(percent*100) / 100
	decimalPart := int(math.Mod(math.Abs(truncated)*100, 100))

	bar := ""

	if decimalPart > 55 {
		restBars--
		bar = strings.Repeat("█", int(math.Ceil(fullBars))) + "▒" + strings.Repeat("░", int(math.Ceil(restBars)))
	} else {
		bar = strings.Repeat("█", int(math.Ceil(fullBars))) + strings.Repeat("░", int(math.Ceil(restBars)))
	}

	return bar
}

// convertMinutesToHours converts minutes to hours.
func convertMinutesToHours(minutes int) (int, int) {
	hours := minutes / 60
	remainingMinutes := minutes % 60
	return hours, remainingMinutes
}

func (mw loggingMiddleware) Wakatime(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		mw.logger.Info(
			"service invocation",
			"method", "Wakatime",
			"result.err", err,
			"took", (time.Since(begin) / 1e6).String(),
		)
	}(time.Now())
	return mw.next.Wakatime(ctx)
}
