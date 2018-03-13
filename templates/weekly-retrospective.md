**Weekly retrospective**

What 3 things going well?
{{ range .GoingWell }}
- {{ . -}}
{{ end }}

What 3 things need improvement?
{{ range .NeedsImprovement }}
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
    - {{ .Title }}
    {{- range .WeeklyGoals }}
        - week {{ .Week }}: {{ .Title }} ({{ .Created }}) {{ .Status -}}
    {{ end -}}
{{ end }}

- How did sprints go:
{{- range .MonthlySprints }}
    - {{ .Title }}
    {{- range .WeeklyGoals }}
        - week {{ .Week }}: {{ .Title }} ({{ .Created }}) {{ .Status -}}
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


----


for reference, completed:
{{range .DoneByDay}}
- {{ .Date }}
{{- range .Done }}
    - {{ .Title }}
{{- end -}}
{{end}}
