package models
import "time"

type Event struct {
  ID          int64      `json:"id"`
  Name        string     `json:"name"`
  Description string     `json:"description"`
  DateTime    time.Time  `json:"date_time"`
  Location    string     `json:"location"`
}
