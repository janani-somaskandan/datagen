package config

import(
    "time"
)

type UserSegment struct {
    Number_of_users int
    Activity_ticker_in_seconds int
    Activity_probablity_map map[string]float64
    Event_probablity_map map[string]float64
    Start_Time time.Time
    Event_attributes map[string]string
    User_attributes map[string]string
    Set_attributes bool
    
}
type Configuration struct {  
    Activity_time_in_seconds int
    Real_Time bool
    User_id_prefix string
    User_segments map[string]UserSegment
}