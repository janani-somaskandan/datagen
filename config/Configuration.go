package config

import(
    "time"
)

type Configuration struct {
    Number_of_users int
    Activity_time_in_seconds int
    Activity_ticker_in_seconds int
    Activity_probablity_map map[string]float64
    Event_probablity_map map[string]float64
    Start_Time time.Time
    Real_Time bool
}