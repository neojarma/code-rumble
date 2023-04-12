package entity

type GameSubmissionPayload struct {
	PlayerId          string             `json:"playerId"`
	RoomId            string             `json:"roomId"`
	TestCaseLimit     int                `json:"testCaseLimit"`
	SubmissionPayload *SubmissionPayload `json:"submissionPayload"`
}
