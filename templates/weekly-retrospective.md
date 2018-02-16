<h1>Week of {{.ThisWeekSunday}}</h1>

**Weekly retrospective**
What 3 things going well?
{{- range .GoingWell }}
- {{ . -}}
{{ end }}

What 3 things need improvement?
{{- range .NeedsImprovement }}
- {{ . -}}
{{ end }}


How did my week go?
- what success did i experience?
{{- range .Successes }}
    - {{ . -}}
{{ end }}

- what challenges did i endure?
{{- range .Challenges }}
    - {{ . -}}
{{ end }}

What did i learn this week?
- about myself?
{{- range .LearnAboutMyself }}
    - {{ . -}}
{{ end }}

- about others?
{{- range .LearnAboutOthers }}
    - {{ . -}}
{{ end }}

Compared to outcomes planned for the week,
- what actually was or wasn’t done:
{{- range .MonthlyGoals }}
    - {{ .Title }} ({{ .Created }})
    {{- range .WeeklyGoals }}
        - week {{ .Week }}: {{ .Title -}} ({{ .Created }}) {{ .Status -}}
    {{ end -}}
{{ end }}

- How did sprints go:
{{- range .MonthlySprints }}
    - {{ .Title }} ({{ .Created }})
    {{- range .WeeklyGoals }}
        - week {{ .Week }}: {{ .Title -}} ({{ .Created }}) {{ .Status -}}
    {{ end -}}
{{ end }}

- what I did to create that particular outcome or situation?
{{- range .WhatIDidToCreateOutcome }}
    - {{ . -}}
{{ end }}

- What do i plan to do differently or the same next week?
{{- range .WhatIPlanToDoDifferently }}
    - {{ . -}}
{{ end }}

Who did i interact with?
- anyone i need to update?
- thank?
- ask a question?
- share feedback with?

[ ] review Yearly goals https://paper.dropbox.com/doc/Yearly-goals-2017-odfgfBzyJr3JFEixro4t0
[ ] review work: task paper
[ ] review timeular (where was time spent?)
[ ] review toggl (where was time spent?)
[ ] review Google DB calendar (on phone) (what happened?)
[ ] review personal calendar (on phone / in Fantastical) (what happened?)
[ ] review Kanban done below (what was done?)
[ ] review goals for week above
[ ] review sprints/habits for week above
[ ] review retrospective notes in Bear

[ ] end here for monthly review, start up monthly-retrospective.md file

[ ] Set new weekly goals in trello list
[ ] Update sprints/habits with priorities/focus for the upcoming week (task paper)
[ ] move work backlog forward to this week backlog (task paper)
[ ] update planning 2018 with `wgoals`
[ ] move weekly review card to done, then `trello-workflow-cli wc`

----


starting {{.NowHHMM}} -


for reference, completed:
{{range .DoneByDay}}
- {{ .Date }}
{{- range .Done }}
    - {{ .Title }}
{{- end -}}
{{end}}
