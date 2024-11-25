// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package queries

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type EventType string

const (
	EventTypeGOAL    EventType = "GOAL"
	EventTypeOWNGOAL EventType = "OWN_GOAL"
	EventTypeYELLOW  EventType = "YELLOW"
	EventTypeRED     EventType = "RED"
	EventTypeSUB     EventType = "SUB"
)

func (e *EventType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = EventType(s)
	case string:
		*e = EventType(s)
	default:
		return fmt.Errorf("unsupported scan type for EventType: %T", src)
	}
	return nil
}

type NullEventType struct {
	EventType EventType
	Valid     bool // Valid is true if EventType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullEventType) Scan(value interface{}) error {
	if value == nil {
		ns.EventType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.EventType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullEventType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.EventType), nil
}

type PlayerPosition string

const (
	PlayerPositionGK  PlayerPosition = "GK"
	PlayerPositionDEF PlayerPosition = "DEF"
	PlayerPositionMFD PlayerPosition = "MFD"
	PlayerPositionFWD PlayerPosition = "FWD"
	PlayerPositionSUB PlayerPosition = "SUB"
)

func (e *PlayerPosition) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PlayerPosition(s)
	case string:
		*e = PlayerPosition(s)
	default:
		return fmt.Errorf("unsupported scan type for PlayerPosition: %T", src)
	}
	return nil
}

type NullPlayerPosition struct {
	PlayerPosition PlayerPosition
	Valid          bool // Valid is true if PlayerPosition is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPlayerPosition) Scan(value interface{}) error {
	if value == nil {
		ns.PlayerPosition, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PlayerPosition.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPlayerPosition) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PlayerPosition), nil
}

type UserRole string

const (
	UserRoleUSER  UserRole = "USER"
	UserRoleADMIN UserRole = "ADMIN"
)

func (e *UserRole) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UserRole(s)
	case string:
		*e = UserRole(s)
	default:
		return fmt.Errorf("unsupported scan type for UserRole: %T", src)
	}
	return nil
}

type NullUserRole struct {
	UserRole UserRole
	Valid    bool // Valid is true if UserRole is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullUserRole) Scan(value interface{}) error {
	if value == nil {
		ns.UserRole, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.UserRole.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullUserRole) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.UserRole), nil
}

type Club struct {
	ID        string
	Name      string
	ShortName string
	Stadium   pgtype.Text
	Logo      string
	Est       int32
}

type ClubPlayer struct {
	ClubID   string
	PlayerID int32
	Season   string
	No       int16
}

type FantasyPlayer struct {
	ClubID   pgtype.Text
	PlayerID pgtype.Int4
	Cost     int32
	Points   pgtype.Int4
	Rating   pgtype.Int4
}

type Lineup struct {
	ID            int32
	ClubID        string
	Possession    pgtype.Numeric
	ShotsOnTarget int16
	Shots         int16
	Touches       int16
	Passes        int16
	Tackles       int16
	Clearances    int16
	Corners       int16
	Offsides      int16
	FoulsConceded int16
}

type LineupEvent struct {
	ID        int32
	LineupID  int32
	PlayerId1 pgtype.Int4
	PlayerId2 pgtype.Int4
	Event     EventType
	Minutes   int16
	Extra     pgtype.Int2
	AfterHalf bool
}

type LineupPlayer struct {
	LineupID   int32
	PlayerID   int32
	PositionNo int16
	Position   PlayerPosition
}

type Match struct {
	ID           int32
	HomeLineupID int32
	AwayLineupID int32
	Season       string
	Week         int16
	Location     string
	StartAt      pgtype.Timestamp
	IsFinished   bool
}

type Player struct {
	ID          int32
	Firstname   string
	Lastname    string
	Dob         pgtype.Timestamp
	Height      int16
	Nationality string
	Position    PlayerPosition
	Image       pgtype.Text
}

type Session struct {
	Token     string
	Username  string
	ExpiresAt pgtype.Timestamp
	CreatedAt pgtype.Timestamp
}

type User struct {
	Username     string
	PasswordHash string
	Firstname    string
	Lastname     string
	Role         UserRole
}
