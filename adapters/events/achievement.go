package events

type AchievementDefinition struct {
	Code      string
	Name      string
	Group     string
	Position  int
	Condition func(PurchaseStats) bool
}

type PurchaseStats struct {
	TotalPurchases int
}

type AchievementUnlockedPayload struct {
	AchievementName string `json:"achievement_name"`
	User            User   `json:"user"`
}
