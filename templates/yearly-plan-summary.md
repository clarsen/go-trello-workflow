# {{ .Year }}
## Outcomes for the year
( see paper doc https://paper.dropbox.com/doc/Yearly-Goals-odfgfBzyJr3JFEixro4t0 )

{{ range .MonthlySummaries }}
## {{ .Month | formatMonthAsString }}

### Outcomes
{{- range .MonthlyGoals }}
#### {{ .Title }} ({{ .Created }})
  {{- range .WeeklyGoals }}
- week {{ .Week }}: {{ .Title }} ({{ .Created }}) {{ .Status -}}
  {{ end }}
{{ end }}

### Events
- list of events for {{ .Month | formatMonthAsString }}

### Sprints
{{- range .MonthlySprints }}
- {{ .Title }} ({{ .Created }})
  {{- range .WeeklyGoals }}
    - week {{ .Week }}: {{ .Title }} ({{ .Created }}) {{ .Status -}}
  {{ end -}}
{{ end }}


{{ end }}
