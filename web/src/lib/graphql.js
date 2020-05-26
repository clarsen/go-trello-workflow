import gql from 'graphql-tag'
import { Mutation } from 'react-apollo'

const fragments = {
  task: gql`
    fragment TaskWhole on Task {
      id
      title
      createdDate
      url
      due
      list {
        board
        list
      }
      period
      dateLastActivity
    }
  `,
  monthlyGoal: gql`
    fragment MonthlyGoalWhole on MonthlyGoal {
      idCard
      title
      weeklyGoals {
        idCard
        idCheckitem
        title
        week
        tasks {
          id
          title
        }
        done
        status
      }

    }
  `

}

export const TaskQuery = gql`
  query tasks($inBoardList: BoardListInput) {
    tasks(inBoardList: $inBoardList) {
      ...TaskWhole
    }
  }
  ${fragments.task}
`

export const MonthlyGoalsQuery = gql`
  query monthlyGoals {
    monthlyGoals {
      ...MonthlyGoalWhole
    }
  }
  ${fragments.monthlyGoal}
`

const PrepareWeeklyReviewQuery = gql`
  mutation prepareWeeklyReview($year: Int, $week: Int) {
    prepareWeeklyReview(year: $year, week: $week) {
      message
      ok
    }
  }
`

export const PrepareWeeklyReview = ({ render }) => (
  <Mutation
    mutation={PrepareWeeklyReviewQuery}
  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)

const finishWeeklyReviewQuery = gql`
  mutation finishWeeklyReview($year: Int, $week: Int) {
    finishWeeklyReview(year: $year, week: $week) {
      message
      ok
    }
  }
`

export const FinishWeeklyReview = ({ render }) => (
  <Mutation
    mutation={finishWeeklyReviewQuery}
  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)

const setDueDateQuery = gql`
  mutation setDueDate($taskId: String!, $due: Timestamp!) {
    setDueDate(taskID: $taskId, due: $due) {
      ...TaskWhole
    }
  }
  ${fragments.task}
`

export const SetDueDate = ({ render }) => (
  <Mutation
    mutation={setDueDateQuery}
  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)

const setDoneQuery = gql`
  mutation setDone($taskId: String!, $done: Boolean!, $titleComment: String, $nextDue: Timestamp) {
    setDone(taskID: $taskId, done: $done, titleComment: $titleComment, nextDue: $nextDue) {
      ...TaskWhole
    }
  }
  ${fragments.task}
`

export const SetDone = ({ render }) => (
  <Mutation
    mutation={setDoneQuery}
  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)

const setGoalDoneQuery = gql`
  mutation setGoalDone($taskId: String!, $checkitemID: String!, $done: Boolean!, $status: String) {
    setGoalDone(taskID: $taskId, checkitemID: $checkitemID, done: $done, status: $status) {
      ...MonthlyGoalWhole
    }
  }
  ${fragments.monthlyGoal}
`

export const SetGoalDone = ({ render }) => (
  <Mutation
    mutation={setGoalDoneQuery}
    update={(cache, { data: { setGoalDone } }) => {
      console.log('mutation update got setGoalDone', setGoalDone)

      const query = MonthlyGoalsQuery
      const { monthlyGoals } = cache.readQuery({ query })
      console.log('currently monthlyGoals', monthlyGoals)

      cache.writeQuery({
        query,
        data: { monthlyGoals: setGoalDone }
      })
    }}  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)

const moveTaskToListQuery = gql`
  mutation moveTaskToList($taskID: String!, $list: BoardListInput!) {
    moveTaskToList(taskID: $taskID, list: $list) {
      ...TaskWhole
    }
  }
  ${fragments.task}
`

export const MoveTaskToList = ({ render }) => (
  <Mutation
    mutation={moveTaskToListQuery}
  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)


export const WeeklyVisualizationQuery = gql`
  query weeklyVisualization($year: Int, $week: Int) {
    weeklyVisualization(year: $year, week: $week)
  }
`

export const MonthlyVisualizationQuery = gql`
  query monthlyVisualization($year: Int, $month: Int) {
    monthlyVisualization(year: $year, month: $month)
  }
`

export const ActiveTimerQuery = gql`
  query activeTimer {
    activeTimer {
      id
      title
    }
  }
`

const stopTimerQuery = gql`
  mutation stopTimer($timerID: String!) {
    stopTimer(timerID: $timerID)
  }
`

export const StopTimer = ({ render }) => (
  <Mutation
    mutation={stopTimerQuery}
    update={(cache, { data: { stopTimer } }) => {
      console.log('mutation update got stopTimer', stopTimer)
      cache.writeQuery({
        query: ActiveTimerQuery,
        data: { 
          activeTimer: null,
        }
      })
    }}

  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)

const startTimerQuery = gql`
  mutation startTimer($taskID: String!, $checkitemID: String) {
    startTimer(taskID: $taskID, checkitemID: $checkitemID) {
      id
      title
    }
  }
`

export const StartTimer = ({ render }) => (
  <Mutation
    mutation={startTimerQuery}
    update={(cache, { data: { startTimer } }) => {
      console.log('mutation update got startTimer', startTimer)
      console.log('cache is', cache)
      cache.writeQuery({
        query: ActiveTimerQuery,
        data: { 
          activeTimer: startTimer,
        }
      })
    }}
  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)

const addTaskQuery = gql`
  mutation addTask($title: String!, $board: String, $list: String) {
    addTask(title: $title, board: $board, list: $list) {
      ...TaskWhole
    }
  }
  ${fragments.task}
`

export const AddTask = ({ render }) => (
  <Mutation
    mutation={addTaskQuery}
    update={(cache, { data: { addTask } }) => {
      console.log('mutation update got addTask', addTask)

      const query = TaskQuery
      const { tasks } = cache.readQuery({ query })
      console.log('currently tasks', tasks)

      cache.writeQuery({
        query,
        data: { tasks: tasks.concat([addTask]) }
      })
    }}
  >

    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)

const addWeeklyGoalQuery = gql`
  mutation addWeeklyGoal($taskID: ID!, $title: String!, $week: Int!) {
    addWeeklyGoal(taskID: $taskID, title: $title, week: $week) {
      ...MonthlyGoalWhole
    }
  }
  ${fragments.monthlyGoal}
`

export const AddWeeklyGoal = ({ render }) => (
  <Mutation
    mutation={addWeeklyGoalQuery}
    update={(cache, { data: { addWeeklyGoal } }) => {
      console.log('mutation update got addWeeklyGoal', addWeeklyGoal)

      const query = MonthlyGoalsQuery
      const { monthlyGoals } = cache.readQuery({ query })
      console.log('currently monthlyGoals', monthlyGoals)

      cache.writeQuery({
        query,
        data: { monthlyGoals: addWeeklyGoal }
      })
    }}  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)

const addMonthlyGoalQuery = gql`
  mutation addMonthlyGoal($title: String!) {
    addMonthlyGoal(title: $title) {
      ...MonthlyGoalWhole
    }
  }
  ${fragments.monthlyGoal}
`

export const AddMonthlyGoal = ({ render }) => (
  <Mutation
    mutation={addMonthlyGoalQuery}
    update={(cache, { data: { addMonthlyGoal } }) => {
      console.log('mutation update got addMonthlyGoal', addMonthlyGoal)

      const query = MonthlyGoalsQuery
      const { monthlyGoals } = cache.readQuery({ query })
      console.log('currently monthlyGoals', monthlyGoals)

      cache.writeQuery({
        query,
        data: { monthlyGoals: addMonthlyGoal }
      })
    }}  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)

const prepareMonthlyReviewQuery = gql`
  mutation prepareMonthlyReview($year: Int, $month: Int) {
    prepareMonthlyReview(year: $year, month: $month) {
      message
      ok
    }
  }
`

export const PrepareMonthlyReview = ({ render }) => (
  <Mutation
    mutation={prepareMonthlyReviewQuery}
  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)

const finishMonthlyReviewQuery = gql`
  mutation finishMonthlyReview($year: Int, $month: Int) {
    finishMonthlyReview(year: $year, month: $month) {
      message
      ok
    }
  }
`
export const FinishMonthlyReview = ({ render }) => (
  <Mutation
    mutation={finishMonthlyReviewQuery}
  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)
