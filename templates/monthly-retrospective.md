**Monthly retrospective**

Compared to outcomes planned for the month, what actually was or wasn’t done and what I did to create that particular outcome or situation?
{{- range .MonthlyGoalReviews }}
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
{{- range .DoDifferently }}
- {{ . -}}
{{ end }}

How did the sprint go?  What to continue, what to change?
{{- range .MonthlySprintReviews }}
- {{ .Title }}
    {{- range .CommentsContinueChange }}
    - {{ . -}}
    {{ end }}
{{ end }}

What candidate goals for upcoming month?
{{- range .CandidateGoals }}
- {{ . -}}
{{ end }}

What candidate sprints for upcoming month?
{{- range .CandidateSprints }}
- {{ . -}}
{{ end }}

Highlights?  Anything to post to results board?
{{- range .Highlights }}
- {{ . -}}
{{ end }}

These could be a checklist in Trello or somewhere else instead.
[ ] monthly goals to outcome bullet items from `wgoals`
[ ] review day one weekly reviews over last month
   [ ] capturing highlights
   [ ] fill in outcome - what I did to create the outcome?
   [ ] fill in what planning to do differently/same going forward
   [ ] review what going/not going well - weekly reviews - what planning to do diff going forward

Determining goals, sprints for next month
[ ] from personal open loops https://paper.dropbox.com/doc/Personal-Open-loops-8ceQvquIRKThtAbEANVok
[ ] from yearly goal dashboard https://paper.dropbox.com/doc/Yearly-goals-2017-odfgfBzyJr3JFEixro4t0
[ ] from work dashboard https://paper.dropbox.com/doc/Dashboard-MLzBn77pvN7Mfx1RpfYAm
[ ] from what to do differently

wrapping up
[ ] move weekly review to Done,
[ ] run weekly review cleanup
    [ ] `trello-workflow-cli wc`

[ ] close out monthly review (copies monthly sprints/goals to history board)
    [ ] `node index.js --monthly-review MONTH`

[ ] set new/updated monthly and weekly goals in trello
    [ ] Set new monthly goals (in board)
    [ ] Update carryover monthly goals (in board)
    [ ] Archive ( `c`) old monthly goals that don't carry over in Kanban board (they've already been copied to History)
    [ ] Set next weekly goals (in board)
    [ ] Set new monthly goals (in planning doc) `wgoals`
    [ ] Set next weekly goals (in planning doc) `wgoals`

[ ] set new/updated monthly sprints in trello
    [ ] Add monthly sprints in Kanban board
    [ ] Update carry over sprints in Kanban board
    [ ] Archive ( `c`) old monthly sprints that don't carry over in Kanban board (they've already been copied to History)
    [ ] Set new monthly sprints (in planning doc)

[ ] move monthly review card to done
[ ] copy monthly review card to history (last week of month)
[ ] move monthly review card back to inbox (it will be moved automatically back to periodic board)


----


starting {{.NowHHMM}} -
