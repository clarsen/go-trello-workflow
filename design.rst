Summary of Trello, intermediate form that is not tied to Trello.  If Trello were
replaced with a different tool, e.g. taskpaper, the intermediate form and
downstream tools would remain unchanged.


..

  task:
    title: <title of task>
    doneDate: <date which task was done>
    createdDate: <date which task was first created>
    period: <if periodic task, how often>

  weeklygoal:
    title: <title of weekly goal>
    createdDate: <date which goal was first created>
    status: "done" | "not done" | "partial" | free form text
    weekNumber: <ISO week of year>
    year: YYYY

  monthlygoal:
    title: <title of goal>
    createdDate: <date which goal was first created>
    weeklyGoals:
      sequence of weeklygoal

  monthlysprint:
    title: <title of sprint goal>
    createdDate: <date which goal was first created>
    weeklyGoals:
      sequence of weeklygoal

  Summary consists of:
  year: YYYY
  weekNumber: <ISO week of year>
  done:
    sequence of task

  monthlyGoals:
    sequence of monthlygoal

  monthlySprints:
    sequence of monthlysprint
