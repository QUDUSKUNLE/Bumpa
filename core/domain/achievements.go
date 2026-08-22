package domain

type AchievementUnlockedPayload struct {
	AchievementName string `json:"achievement_name"`
	User            User   `json:"user"`
}
