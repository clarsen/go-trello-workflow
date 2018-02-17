package workflow

// WeeklyReview defines the manually input data that goes into weekly review visualization
type WeeklyReview struct {
	GoingWell                []string `yaml:"goingWell"`
	NeedsImprovement         []string `yaml:"needsImprovement"`
	Successes                []string `yaml:"successes"`
	Challenges               []string `yaml:"challenges"`
	LearnAboutMyself         []string `yaml:"learnAboutMyself"`
	LearnAboutOthers         []string `yaml:"learnAboutOthers"`
	WhatIDidToCreateOutcome  []string `yaml:"whatIDidToCreateOutcome"`
	WhatIPlanToDoDifferently []string `yaml:"whatIPlanToDoDifferently"`
}

type MonthlyGoalReview struct {
	Title           string   `yaml:"title"` // should match Title of MonthlYGoalInfo
	Accomplishments []string `yaml:"accomplishments"`
	CreatedBy       []string `yaml:"createdBy"`
}

type MonthlySprintReview struct {
	Title                  string   `yaml:"title"` // should match Title of MonthlYGoalInfo
	CommentsContinueChange []string `yaml:"commentsContinueChange"`
}

// MonthlyReview defines the manually input data that goes into the monthly
// review visualization
type MonthlyReview struct {
	MonthlyGoalReviews   []MonthlyGoalReview   `yaml:"monthlyGoalReviews"`
	MonthlySprintReviews []MonthlySprintReview `yaml:"monthlySprintReviews"`
	DoDifferently        []string              `yaml:"doDifferently"`
	CandidateGoals       []string              `yaml:"candidateGoals"`
	CandidateSprints     []string              `yaml:"candidateSprints"`
	Highlights           []string              `yaml:"highlights"`
}
