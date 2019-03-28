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
    }
  `,
  monthlyGoal: gql`
    fragment MonthlyGoalWhole on MonthlyGoal {
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

export const PrepareWeeklyReviewQuery = gql`
  mutation prepareWeeklyReview($year: Int, $week: Int) {
    prepareWeeklyReview(year: $year, week: $week) {
      message
      ok
    }
  }
`

export const FinishWeeklyReviewQuery = gql`
  mutation finishWeeklyReview($year: Int, $week: Int) {
    finishWeeklyReview(year: $year, week: $week) {
      message
      ok
    }
  }
`

export const SetDueDateQuery = gql`
  mutation setDueDate($taskId: String!, $due: Timestamp!) {
    setDueDate(taskID: $taskId, due: $due) {
      ...TaskWhole
    }
  }
  ${fragments.task}
`

export const SetDoneQuery = gql`
  mutation setDone($taskId: String!, $done: Boolean!, $nextDue: Timestamp) {
    setDone(taskID: $taskId, done: $done, nextDue: $nextDue) {
      ...TaskWhole
    }
  }
  ${fragments.task}
`

export const SetGoalDoneQuery = gql`
  mutation setGoalDone($taskId: String!, $checkitemID: String!, $done: Boolean!, $status: String) {
    setGoalDone(taskID: $taskId, checkitemID: $checkitemID, done: $done, status: $status) {
      ...MonthlyGoalWhole
    }
  }
  ${fragments.monthlyGoal}
`

export const MoveTaskToListQuery = gql`
  mutation moveTaskToList($taskID: String!, $list: BoardListInput!) {
    moveTaskToList(taskID: $taskID, list: $list) {
      ...TaskWhole
    }
  }
  ${fragments.task}
`

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

export const StopTimerQuery = gql`
  mutation stopTimer($timerID: String!) {
    stopTimer(timerID: $timerID)
  }
`

export const StartTimerQuery = gql`
  mutation startTimer($taskID: String!, $checkitemID: String) {
    startTimer(taskID: $taskID, checkitemID: $checkitemID) {
      id
      title
    }
  }
`

export const AddTaskQuery = gql`
  mutation addTask($title: String!, $board: String, $list: String) {
    addTask(title: $title, board: $board, list: $list) {
      ...TaskWhole
    }
  }
  ${fragments.task}
`

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
