package events

type BadgeDefinition struct {
	Code            string
	Name            string
	RequiredRewards int
}

type BadgeUnlockedPayload struct {
	BadgeName string `json:"badge_name"`
	BadgeCode string `json:"badge_code"`
	User      string   `json:"user"`
}
