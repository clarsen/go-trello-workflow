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

The outcomes planned for the week,
{{ range $index, $ := .PerGoalReviews }}
- {{ index .DidToCreateOutcome 0 }}
    - What did I do to create that particular outcome or situation?
{{- range .DidToCreateOutcome }}
        - {{ . }}
{{ end }}
{{- if .KeepDoing }}
    - What do I plan to keep doing?
{{- range .KeepDoing }}
        - {{ . }}
{{ end -}}
{{ end -}}
{{- if .DoDifferently }}
    - What do I plan to stop doing/do differently?
{{- range .DoDifferently }}
        - {{ . }}
{{ end -}}
{{ end }}
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
