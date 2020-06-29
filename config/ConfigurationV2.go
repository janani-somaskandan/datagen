package config

import(
    "time"
)

type CorrelationMatrix struct {
    Events map[string]map[string]float64
    Seed_events []string
}

type EventProbablity struct{
    Correlation_matrix CorrelationMatrix
    Independent_events map[string]float64
}

type UserSegmentV2 struct {
    Number_of_users int
    Activity_ticker_in_seconds int
    Activity_probablity_map map[string]float64
    Event_probablity_map EventProbablity
    Start_Time time.Time
    Event_attributes map[string]string
    User_attributes map[string]string
    Set_attributes bool
    
}
type ConfigurationV2 struct {  
    Output_file_name string
    Activity_time_in_seconds int
    Real_Time bool
    User_id_prefix string
    User_segments map[string]UserSegmentV2
    New_user_poll_time int
    New_user_probablity float64
}