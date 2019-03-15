import React from 'react'
import {
  Button,
  Container,
  Row,
  Col,
  Spinner
} from 'reactstrap'

import { Query, Mutation } from 'react-apollo'
import { adopt } from 'react-adopt'
import gql from 'graphql-tag'
import { withAlert } from 'react-alert'
import moment from 'moment'
import { FaSync } from 'react-icons/fa'
import NavHeader from '../components/NavHeader'
import TaskList from '../components/TaskList'
import GoalList from '../components/GoalList'
import Timer from '../components/Timer'
import MarkdownRenderer from 'react-markdown-renderer'


export const taskQuery = gql`
  query tasks($inBoardList: BoardListInput) {
    tasks(inBoardList: $inBoardList) {
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
  }
`

export const monthlyGoalsQuery = gql`
  query monthlyGoals {
    monthlyGoals {
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
      }
    }
  }
`

const prepareWeeklyReviewQuery = gql`
  mutation prepareWeeklyReview($year: Int, $week: Int) {
    prepareWeeklyReview(year: $year, week: $week) {
      message
      ok
    }
  }
`

const prepareWeeklyReview = ({ render }) => (
  <Mutation
    mutation={prepareWeeklyReviewQuery}
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

const finishWeeklyReview = ({ render }) => (
  <Mutation
    mutation={finishWeeklyReviewQuery}
  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)

const setDueDateQuery = gql`
mutation setDueDate($taskId: String!, $due: Timestamp!) {
  setDueDate(taskID: $taskId, due: $due) {
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
}
`

const setDueDate = ({ render }) => (
  <Mutation
    mutation={setDueDateQuery}
  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)

const setDoneQuery = gql`
mutation setDone($taskId: String!, $done: Boolean!, $nextDue: Timestamp) {
  setDone(taskID: $taskId, done: $done, nextDue: $nextDue) {
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
}
`

const setDone = ({ render }) => (
  <Mutation
    mutation={setDoneQuery}
  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)

const moveTaskToListQuery = gql`
mutation moveTaskToList($taskID: String!, $list: BoardListInput!) {
  moveTaskToList(taskID: $taskID, list: $list) {
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
}
`

const moveTaskToList = ({ render }) => (
  <Mutation
    mutation={moveTaskToListQuery}
  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)

const weeklyVisualizationQuery = gql`
  query weeklyVisualization($year: Int, $week: Int) {
    weeklyVisualization(year: $year, week: $week)
  }
`

const activeTimerQuery = gql`
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

const stopTimer = ({ render }) => (
  <Mutation
    mutation={stopTimerQuery}
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

const startTimer = ({ render }) => (
  <Mutation
    mutation={startTimerQuery}
  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)



/* eslint-disable react/display-name */
/* eslint-disable react/prop-types */

const QueryContainer = adopt({
  queryAll: ({ render }) => (
    <Query query={taskQuery} ssr={false} variables={{ }}>
      {render}
    </Query>
  ),
  queryAllGoals: ({ render }) => (
    <Query query={monthlyGoalsQuery} ssr={false} variables={{ }}>
      {render}
    </Query>
  ),
  weeklyVisualizationQuery: ({ render, year, week }) => (
    <Query query={weeklyVisualizationQuery} ssr={false} variables={{ year, week }}>
      {render}
    </Query>
  ),
  queryTimer: ({ render }) => (
    <Query query={activeTimerQuery} ssr={false} variables={{ }}>
      {render}
    </Query>
  ),
  prepareWeeklyReview,
  finishWeeklyReview,
  setDueDate,
  setDone,
  moveTaskToList,
  stopTimer,
  startTimer,
})

class IndexPage extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      data: null
    }
  }
  render () {
    let { alert } = this.props
    let now = moment()
    let nowGrace = moment().subtract(3,'days')

    return (
      <QueryContainer year={now.year()} week={now.isoWeek()}>
        {({
          queryAll: { loading: loadingAll, data: allTasks, error: queryAllError, refetch: allRefetch },
          queryAllGoals: { loading: loadingAllGoals, data: allGoals, error: queryAllGoalsError, refetch: allGoalsRefetch },
          weeklyVisualizationQuery: { loading : weeklyLoading, data: weeklyVisualizationData, error: weeklyError, refetch: weeklyVisualizationRefetch },
          queryTimer: { loading: loadingTimer, data: timerData, error: timerError, refetch: timerRefetch },
          prepareWeeklyReview,
          finishWeeklyReview,
          setDueDate,
          setDone,
          moveTaskToList,
          stopTimer,
          startTimer,
        }) =>
          <React.Fragment>
            <NavHeader/>
            <Button color='primary' size='sm' onClick={() => {
              prepareWeeklyReview
                .mutation({
                  variables: {
                    year: now.year(),
                    week: now.isoWeek(),
                  }
                })
                .then(({ data }) => alert.show(data.prepareWeeklyReview.message))
            }}>
                Prepare weekly review for {now.year()}-{now.isoWeek()}
            </Button>{' '}
            {nowGrace.isoWeek() !== now.isoWeek() &&
              <Button color='primary' size='sm' onClick={() => {
                prepareWeeklyReview
                  .mutation({
                    variables: {
                      year: nowGrace.year(),
                      week: nowGrace.isoWeek(),
                    }
                  })
                  .then(({ data }) => alert.show(data.prepareWeeklyReview.message))
              }}>
                  Prepare weekly review for {nowGrace.year()}-{nowGrace.isoWeek()}
              </Button>
            }{' '}
            <Button color='primary' size='sm' onClick={() => {
              finishWeeklyReview
                .mutation({
                  variables: {
                    year: now.year(),
                    week: now.isoWeek(),
                  }
                })
                .then(({ data }) => alert.show(data.finishWeeklyReview.message))
            }}>
                Finish weekly review for {now.year()}-{now.isoWeek()}
            </Button>{' '}
            <FaSync size={25} onClick={() => {
              allRefetch()
              allGoalsRefetch()
            }}/>

            <Container>
              <Row>
                <Col lg={6}>
                  <div className="listTitle">Goals</div>
                  {loadingAllGoals && <Spinner color="primary" />}
                  {!loadingAllGoals && console.log('got data', allGoals)}
                  {queryAllGoalsError && <div>Goals: {queryAllError.message}</div>}
                  {!loadingAllGoals && !queryAllGoalsError && <GoalList startTimer={startTimer} timerRefetch={timerRefetch} goals={allGoals.monthlyGoals}/>}
                </Col>
                <Col lg={6}>
                  {loadingTimer && <Spinner color="primary" />}
                  {!loadingTimer && console.log('got data', timerData)}
                  {timerError && <div>Timer: {timerError.message}</div>}
                  {!loadingTimer && !timerError && <Timer stopTimer={stopTimer} timerRefetch={timerRefetch} activeTimer={timerData.activeTimer} />}
                </Col>
              </Row>
              <Row>
                <Col lg={6}>
                  <div className="listTitle">Today</div>
                  {loadingAll && <Spinner color="primary" />}
                  {!loadingAll && console.log('got data', allTasks)}
                  {queryAllError && <div>Tasks: {queryAllError.message}</div>}
                  {(!loadingAll && !queryAllError) && <TaskList listFilter={['Today']} setDueDate={setDueDate} setDone={setDone} moveTaskToList={moveTaskToList} startTimer={startTimer} timerRefetch={timerRefetch} tasks={allTasks.tasks}/>}
                </Col>
                <Col lg={6}>
                  <div className="listTitle">Waiting on...</div>
                  {loadingAll && <Spinner color="primary" />}
                  {(!loadingAll && !queryAllError) && <TaskList listFilter={['Waiting on']} setDueDate={setDueDate} setDone={setDone} moveTaskToList={moveTaskToList}  startTimer={startTimer} timerRefetch={timerRefetch} tasks={allTasks.tasks}/>}
                </Col>
              </Row>
              <Row>
                <Col lg={6}>
                  <div className="listTitle">Backlog</div>
                  {loadingAll && <Spinner color="primary" />}
                  {(!loadingAll && !queryAllError) && <TaskList listFilter={['Backlog (Personal)']} setDueDate={setDueDate} setDone={setDone} moveTaskToList={moveTaskToList}  startTimer={startTimer} timerRefetch={timerRefetch} tasks={allTasks.tasks}/>}
                </Col>
                <Col lg={6}>
                  <Row>
                    <div className="listTitle">Periodic</div>
                    <div className="listSubGroupTitle">Often</div>
                    {loadingAll && <Spinner color="primary" />}
                    {(!loadingAll && !queryAllError) && <TaskList isPeriodic listFilter={['Often']} setDueDate={setDueDate} setDone={setDone} moveTaskToList={moveTaskToList}  startTimer={startTimer} timerRefetch={timerRefetch} tasks={allTasks.tasks}/>}
                  </Row>
                  <Row>
                    <div className="listSubGroupTitle">Weekly</div>
                    {loadingAll && <Spinner color="primary" />}
                    {(!loadingAll && !queryAllError) && <TaskList noHeader isPeriodic listFilter={['Weekly']} setDueDate={setDueDate} setDone={setDone} startTimer={startTimer} timerRefetch={timerRefetch} moveTaskToList={moveTaskToList}  tasks={allTasks.tasks}/>}
                  </Row>
                  <Row>
                    <div className="listSubGroupTitle">Bi-weekly to monthly</div>
                    {loadingAll && <Spinner color="primary" />}
                    {(!loadingAll && !queryAllError) && <TaskList noHeader isPeriodic listFilter={['Bi-weekly to monthly']} setDueDate={setDueDate} setDone={setDone} startTimer={startTimer} timerRefetch={timerRefetch} moveTaskToList={moveTaskToList}  tasks={allTasks.tasks}/>}
                  </Row>
                  <Row>
                    <div className="listSubGroupTitle">Quarterly to Yearly</div>
                    {loadingAll && <Spinner color="primary" />}
                    {(!loadingAll && !queryAllError) && <TaskList noHeader isPeriodic listFilter={['Quarterly to Yearly']} setDueDate={setDueDate} setDone={setDone} startTimer={startTimer} timerRefetch={timerRefetch} moveTaskToList={moveTaskToList}  tasks={allTasks.tasks}/>}
                  </Row>
                </Col>
              </Row>
              <Row>
                <Col>
                  {weeklyLoading && <Spinner color="primary" />}
                  {(!weeklyLoading && !weeklyError) &&
                    <div>
                      <FaSync size={25} onClick={() => weeklyVisualizationRefetch()} />
                      <MarkdownRenderer className='weeklyReview' markdown={weeklyVisualizationData.weeklyVisualization} />
                    </div>
                  }
                  {weeklyError && <div>Weekly review: {weeklyError.message}</div>}
                </Col>
              </Row>
            </Container>
            <style jsx global>{`
              .listSubGroupTitle {
                background: #999;
              }
              .listTitle {
                background: #bbb;
                width: 100%;
                color: #fff;
              }
              .weeklyReview {
                border: 1px solid #fff;
              }
            `}</style>

          </React.Fragment>
        }
      </QueryContainer>
    )
  }
}

export default withAlert()(IndexPage)
