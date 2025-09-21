package schemas

type OnlineTestSchema struct {
	Title     string           `json:"title" binding:"required"`
	TestID    string           `json:"test_id" binding:"required"`
	Duration  int              `json:"duration" binding:"required"`
	Questions []QuestionSchema `json:"questions" binding:"required,dive"`
}

type QuestionSchema struct {
	Type       string   `json:"type" binding:"required,oneof=mcq tof essay"`
	Text       string   `json:"text" binding:"required"`
	Answer     string   `json:"answer" binding:"required"`
	Options    []string `json:"options"`    // For MCQ
	Statements []string `json:"statements"` // For TOF
}

type AllOnlineTestsSchema struct {
	Tests []OnlineTestSummarySchema `json:"tests"`
}

type OnlineTestSummarySchema struct {
	ID       uint   `json:"id"`
	Title    string `json:"title" binding:"required"`
	TestID   string `json:"test_id" binding:"required"`
	Duration int    `json:"duration" binding:"required"`
}
