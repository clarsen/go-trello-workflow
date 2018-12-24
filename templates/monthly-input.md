**Monthly retrospective (raw input for review)**

Covering weeks {{ .WeeksOfYear }}

# Highlights?  Anything to post to results board?
## What was a success?
{{ range .Successes }}
- Week {{ .Week -}}
    {{ range .Content }}
    - {{ . -}}
    {{ end }}
{{ end }}

# What do I plan to continue doing next month?
## What was going well
{{ range .GoingWell }}
- Week {{ .Week -}}
    {{ range .Content }}
    - {{ . -}}
    {{ end }}
{{ end }}

# What do I plan to do differently next month?
## What needed improvement
{{ range .NeedsImprovement }}
- Week {{ .Week -}}
    {{ range .Content }}
    - {{ . -}}
    {{ end }}
{{ end }}
## What was challenging
{{ range .Challenges }}
- Week {{ .Week -}}
    {{ range .Content }}
    - {{ . -}}
    {{ end }}
{{ end }}

# For each goal
{{ range .MonthlyGoalSummaries }}
## {{ .Goal }}
  - Did to create outcome
    {{ range .DidToCreateOutcome }}
    - Week {{ .Week -}}
        {{ range .Content }}
        - {{ . -}}
        {{ end }}
    {{ end }}
  - Keep doing
    {{ range .KeepDoing }}
    - Week {{ .Week -}}
        {{ range .Content }}
        - {{ . -}}
        {{ end }}
    {{ end }}
  - Do differently
    {{ range .DoDifferently }}
    - Week {{ .Week -}}
        {{ range .Content }}
        - {{ . -}}
        {{ end }}
    {{ end }}

{{ end}}
