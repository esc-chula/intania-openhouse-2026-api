package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/esc-chula/intania-openhouse-2026-api/pkg/config"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/database"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

var (
	urlFlag         = flag.String("url", "http://localhost:8000", "Base URL of the API")
	authFlag        = flag.String("auth", "", "Authorization header value (e.g., 'Bearer token')")
	rateFlag        = flag.Int("rate", 50, "Request rate per second")
	durationFlag    = flag.Duration("duration", 5*time.Second, "Duration of the test")
	connectionsFlag = flag.Int("connections", 10000, "Number of concurrent connections")
	cfgFile         = flag.String("config", "", "Path to config file")
)

func main() {
	flag.Parse()

	if *authFlag == "" {
		log.Fatal("Error: -auth flag is required")
	}

	cfg, err := config.InitConfig(*cfgFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db := database.NewPostgresDB(cfg.Database())
	ctx := context.Background()

	seededData := seedData(ctx, db)
	if err := ensureUserCreated(*urlFlag, *authFlag); err != nil {
		log.Fatalf("failed to ensure user created: %v", err)
	}

	rate := vegeta.Rate{Freq: *rateFlag, Per: time.Second}
	duration := *durationFlag

	targeter := NewCustomTargeter(*urlFlag, *authFlag, seededData)
	attacker := vegeta.NewAttacker(
		vegeta.Connections(*connectionsFlag),
	)

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Load Test") {
		metrics.Add(res)
	}
	metrics.Close()

	reporter := vegeta.NewTextReporter(&metrics)
	reporter.Report(os.Stdout)
}

func NewCustomTargeter(baseURL string, auth string, data *SeedData) vegeta.Targeter {
	headers := http.Header{}
	headers.Add("Authorization", auth)
	headers.Add("Content-Type", "application/json")

	totalActivities := len(data.Activities)
	totalWorkshops := len(data.Workshops)
	totalBooths := len(data.Booths)

	redemptionList := []string{"department", "club", "exhibition"}

	return func(t *vegeta.Target) error {
		if t == nil {
			return vegeta.ErrNilTarget
		}

		randWorkshop := data.Workshops[rand.Intn(totalWorkshops)]
		randBooth := data.Booths[rand.Intn(totalBooths)]
		randActivity := data.Booths[rand.Intn(totalActivities)]
		randRedemption := redemptionList[rand.Intn(3)]

		// Randomly select an endpoint to attack
		endpoints := []struct {
			Method string
			URL    string
			Body   []byte
		}{
			// // CreateUser
			// {"POST", "/users", []byte(`{
			// "first_name": "John",
			// "last_name": "Doe",
			// "gender": "male",
			// "phone_number": "0123456789",
			// "participant_type": "student",
			// "attendance_dates": ["2026-03-28"],
			// "interested_activities": [],
			// "discovery_channel": [],
			// "is_from_bangkok": true,
			// "origin_location": "phra_nakhon",
			// "transport_mode": "personal_car",
			// "student_extra_attributes": {
			//   "education_level": "high school",
			//   "emergency_contact": "0123456789",
			//   "interested_major": "",
			//   "province": "",
			//   "school_name": "",
			//   "study_plan": "",
			//   "tcas_rank": 0
			// }`)},

			// GetUser
			{"GET", "/users/me", nil},

			// GetActivity
			{"GET", fmt.Sprintf("/activities/%d", randActivity.ID), nil},
			// ListActivity
			{"GET", "/activities", nil},

			// GetWorkshop
			{"GET", fmt.Sprintf("/workshops/%d", randWorkshop.ID), nil},
			// ListWorkshop
			{"GET", "/workshops", nil},

			// GetMyBookings
			{"GET", "/users/me/bookings", nil},
			// BookWorkshop
			{"POST", fmt.Sprintf("/workshops/%d/book", randWorkshop.ID), nil},
			// CancelBooking
			{"DELETE", fmt.Sprintf("/workshops/%d/book", randWorkshop.ID), nil},

			// CheckIn
			{"POST", "/check-in", fmt.Appendf(nil, `{"code":"%s"}`, "B-"+randBooth.CheckInCode)},

			// GetUserStamps
			{"GET", "/users/me/stamps", nil},
			// GetRedemptionStatus
			{"GET", "/users/me/redemption-status", nil},
			// RedeemStamps
			{"POST", "/stamps/redemptions" + "?category=" + randRedemption, nil},
		}

		selected := endpoints[rand.Intn(len(endpoints))]

		t.Method = selected.Method
		t.URL = baseURL + selected.URL
		t.Header = headers
		t.Body = selected.Body

		return nil
	}
}
