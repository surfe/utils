package utils

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDaysUntil(t *testing.T) {
	t.Skip()
	t.Parallel()

	tests := []struct {
		name         string
		currentTime  time.Time
		targetTime   time.Time
		expectedDays int
	}{
		{
			name:         "FutureTime",
			targetTime:   time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, time.Now().Hour()-1, 0, 0, 0, time.UTC),
			expectedDays: 1,
		},
		{
			name:         "PastTime",
			targetTime:   time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, time.Now().Hour()-1, 0, 0, 0, time.UTC),
			expectedDays: -1,
		},
		{
			name:         "SameTime",
			targetTime:   time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour()-1, 0, 0, 0, time.UTC),
			expectedDays: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			days := DaysUntil(test.targetTime)
			if days != test.expectedDays {
				t.Errorf("Expected %d days until %v, but got %d", test.expectedDays, test.targetTime, days)
			}
		})
	}
}

func TestUTCTimestampToDateOrDefault(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "Correct UTC time and date",
			args: "1995-09-15T15:04:05.000Z",
			want: "15/09/1995",
		},
		{
			name: "Correct time and date, but in the wrong format 1",
			args: "1995/15/09",
			want: "1995/15/09",
		},
		{
			name: "Correct time and date, but in the wrong format 2",
			args: "09/15/1995",
			want: "09/15/1995",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := UTCTimestampToDateOrDefault(tt.args)
			require.Equal(t, tt.want, result)
		})
	}
}

func TestGetTimezone(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		timezone string
		expected string
	}{
		{
			name:     "Valid timezone",
			timezone: "America/New_York",
			expected: "America/New_York",
		},
		{
			name:     "Invalid timezone",
			timezone: "Invalid/Timezone",
			expected: "UTC",
		},
		{
			name:     "Return timezone",
			timezone: "Asia/Tokyo", // Replace with any valid timezone
			expected: "Asia/Tokyo",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			result := GetTimezone(ctx, tc.timezone)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestGetTimeFromLinkedInSentTimeLabelAndUserTimeZone(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC().In(time.UTC)
	testCases := []struct {
		name          string
		sentTimeLabel string
		expectedTime  time.Time
	}{
		{
			name:          "Linkedin invite sent yesterday",
			sentTimeLabel: " yesterday",
			expectedTime:  now.AddDate(0, 0, -1),
		},
		{
			name:          "Linkedin invite sent 2 months ago",
			sentTimeLabel: "2 months",
			expectedTime:  now.AddDate(0, -2, 0),
		},
		{
			name:          "Linkedin invite sent 3 weeks ago",
			sentTimeLabel: "3 weeks",
			expectedTime:  now.AddDate(0, 0, -7*3),
		},
		{
			name:          "Linkedin invite sent 4 days ago",
			sentTimeLabel: "4 days",
			expectedTime:  now.AddDate(0, 0, -4),
		},
		{
			name:          "Linkedin invite sent today",
			sentTimeLabel: "today",
			expectedTime:  now.AddDate(0, 0, 0),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			result := GetTimeFromLinkedInSentTimeLabelAndUserTimeZone(ctx, tc.sentTimeLabel, "UTC")
			y1, m1, d1 := tc.expectedTime.Date()
			y2, m2, d2 := result.Date()
			require.Equal(t, y1, y2)
			require.Equal(t, m1, m2)
			require.Equal(t, d1, d2)
		})
	}
}

func TestParseISODuration(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		duration      string
		expected      time.Duration
		expectError   bool
		errorContains string
	}{
		{
			name:        "Only days",
			duration:    "P3D",
			expected:    3 * 24 * time.Hour,
			expectError: false,
		},
		{
			name:        "Only hours",
			duration:    "PT5H",
			expected:    5 * time.Hour,
			expectError: false,
		},
		{
			name:        "Only minutes",
			duration:    "PT30M",
			expected:    30 * time.Minute,
			expectError: false,
		},
		{
			name:        "Only seconds",
			duration:    "PT45S",
			expected:    45 * time.Second,
			expectError: false,
		},
		{
			name:        "Only years",
			duration:    "P2Y",
			expected:    2 * 365 * 24 * time.Hour,
			expectError: false,
		},
		{
			name:        "Only months",
			duration:    "P6M",
			expected:    6 * 30 * 24 * time.Hour,
			expectError: false,
		},
		{
			name:        "Only weeks",
			duration:    "P4W",
			expected:    4 * 7 * 24 * time.Hour,
			expectError: false,
		},
		{
			name:        "Mixed date and time",
			duration:    "P2Y3M4DT5H6M7S",
			expected:    2*365*24*time.Hour + 3*30*24*time.Hour + 4*24*time.Hour + 5*time.Hour + 6*time.Minute + 7*time.Second,
			expectError: false,
		},
		{
			name:        "Mixed date only",
			duration:    "P1Y2M3W4D",
			expected:    1*365*24*time.Hour + 2*30*24*time.Hour + 3*7*24*time.Hour + 4*24*time.Hour,
			expectError: false,
		},
		{
			name:        "Mixed time only",
			duration:    "PT1H2M3S",
			expected:    1*time.Hour + 2*time.Minute + 3*time.Second,
			expectError: false,
		},
		{
			name:          "Missing P prefix",
			duration:      "1Y2M",
			expectError:   true,
			errorContains: "must start with P",
		},
		{
			name:          "Invalid date format",
			duration:      "PXY",
			expectError:   true,
			errorContains: "invalid date duration format",
		},
		{
			name:          "Invalid time format",
			duration:      "PTXH",
			expectError:   true,
			errorContains: "invalid time duration format",
		},
		{
			name:        "Empty date part with time",
			duration:    "PT5H",
			expected:    5 * time.Hour,
			expectError: false,
		},
		{
			name:        "Empty time part with date",
			duration:    "P5D",
			expected:    5 * 24 * time.Hour,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := ParseISODuration(tc.duration)

			if tc.expectError {
				require.Error(t, err)

				if tc.errorContains != "" && err != nil {
					require.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, result)
			}
		})
	}
}
