**Monthly retrospective**

Compared to outcomes planned for the month, what actually was or wasn’t done and what I did to create that particular outcome or situation?
{{ range .MonthlyGoalReviews }}
- {{ .Title -}}
    {{ range .Accomplishments }}
    - {{ . -}}
    {{ end }}
    - Created by:
    {{- range .CreatedBy }}
        - {{ . -}}
    {{ end }}
{{ end }}

What do I plan to do differently or the same next month?
{{ range .DoDifferently }}
- {{ . -}}
{{ end }}

How did the sprint go?  What to continue, what to change?
{{ range .MonthlySprintReviews }}
- {{ .Title }}
    {{- range .CommentsContinueChange }}
    - {{ . -}}
    {{ end }}
{{ end }}

What candidate goals for upcoming month?
{{ range .CandidateGoals }}
- {{ . -}}
{{ end }}

What candidate sprints for upcoming month?
{{ range .CandidateSprints }}
- {{ . -}}
{{ end }}

Highlights?  Anything to post to results board?
{{ range .Highlights }}
- {{ . -}}
{{ end }}


----

starting {{.NowHHMM}} -
