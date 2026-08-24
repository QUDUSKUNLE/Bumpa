package events

import "github.com/jackc/pgx/v5/pgtype"

type BadgeDefinition struct {
	Code            string
	Name            string
	RequiredRewards int
}

type BadgeUnlockedPayload struct {
	BadgeName string      `json:"badge_name"`
	BadgeCode string      `json:"badge_code"`
	User      pgtype.UUID `json:"user"`
}
