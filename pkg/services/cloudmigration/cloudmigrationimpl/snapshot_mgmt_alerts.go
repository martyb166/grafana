package cloudmigrationimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/alertmanager/timeinterval"

	"github.com/grafana/grafana/pkg/services/cloudmigration/slicesext"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/ngalert/api/tooling/definitions"
	"github.com/grafana/grafana/pkg/services/user"
)

type timeRange struct {
	StartMinute int
	EndMinute   int
}

type inclusiveRange struct {
	Begin int
	End   int
}

type timeInterval struct {
	Location    *time.Location   `json:"location,omitempty"`
	Times       []timeRange      `json:"times,omitempty"`
	Weekdays    []inclusiveRange `json:"weekdays,omitempty"`
	DaysOfMonth []inclusiveRange `json:"days_of_month,omitempty"`
	Months      []inclusiveRange `json:"months,omitempty"`
	Years       []inclusiveRange `json:"years,omitempty"`
}

func timeIntervalFromModel(interval timeinterval.TimeInterval) timeInterval {
	var location *time.Location
	if interval.Location != nil {
		location = interval.Location.Location
	}

	return timeInterval{
		Times: slicesext.Map(interval.Times, func(v timeinterval.TimeRange) timeRange {
			return timeRange{
				StartMinute: v.StartMinute,
				EndMinute:   v.EndMinute,
			}
		}),
		Weekdays: slicesext.Map(interval.Weekdays, func(v timeinterval.WeekdayRange) inclusiveRange {
			return inclusiveRange{
				Begin: v.Begin,
				End:   v.End,
			}
		}),
		DaysOfMonth: slicesext.Map(interval.DaysOfMonth, func(v timeinterval.DayOfMonthRange) inclusiveRange {
			return inclusiveRange{
				Begin: v.Begin,
				End:   v.End,
			}
		}),
		Months: slicesext.Map(interval.Months, func(v timeinterval.MonthRange) inclusiveRange {
			return inclusiveRange{
				Begin: v.Begin,
				End:   v.End,
			}
		}),
		Years: slicesext.Map(interval.Years, func(v timeinterval.YearRange) inclusiveRange {
			return inclusiveRange{
				Begin: v.Begin,
				End:   v.End,
			}
		}),
		Location: location,
	}
}

type muteTimeInterval struct {
	UID           string         `json:"uid"`
	Name          string         `json:"name"`
	Version       string         `json:"version,omitempty"`
	TimeIntervals []timeInterval `json:"time_intervals"`
}

func muteTimeIntervalFromModel(timeInterval *definitions.MuteTimeInterval) muteTimeInterval {
	return muteTimeInterval{
		UID:           timeInterval.UID,
		Name:          timeInterval.Name,
		TimeIntervals: slicesext.Map(timeInterval.TimeIntervals, timeIntervalFromModel),
		Version:       timeInterval.Version,
	}
}

func (s *Service) getAlertMuteTimings(ctx context.Context, signedInUser *user.SignedInUser) ([]muteTimeInterval, error) {
	if !s.features.IsEnabledGlobally(featuremgmt.FlagOnPremToCloudMigrationsAlerts) {
		return nil, nil
	}

	muteTimings, err := s.ngAlert.Api.MuteTimings.GetMuteTimings(ctx, signedInUser.OrgID)
	if err != nil {
		return nil, fmt.Errorf("fetching ngalert mute timings: %w", err)
	}

	muteTimeIntervals := make([]muteTimeInterval, 0, len(muteTimings))

	for _, muteTiming := range muteTimings {
		muteTimeIntervals = append(muteTimeIntervals, muteTimeIntervalFromModel(&muteTiming))
	}

	return muteTimeIntervals, nil
}
