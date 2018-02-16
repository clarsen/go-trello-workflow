package workflow

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
