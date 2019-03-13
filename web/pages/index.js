import React from 'react'
import {
  Button,
  Container,
  Row,
  Col
} from 'reactstrap'

import { Query, Mutation } from 'react-apollo'
import { adopt } from 'react-adopt'
import gql from 'graphql-tag'
import { withAlert } from 'react-alert'
import moment from 'moment'
import { FaSync } from 'react-icons/fa'
import NavHeader from '../components/NavHeader'
import TaskList from '../components/TaskList'
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

const weeklyVisualizationQuery = gql`
  query weeklyVisualization($year: Int, $week: Int) {
    weeklyVisualization(year: $year, week: $week)
  }
`

/* eslint-disable react/display-name */
/* eslint-disable react/prop-types */

const QueryContainer = adopt({
  queryAll: ({ render }) => (
    <Query query={taskQuery} ssr={false} variables={{ }}>
      {render}
    </Query>
  ),
  weeklyVisualizationQuery: ({ render, year, week }) => (
    <Query query={weeklyVisualizationQuery} ssr={false} variables={{ year, week }}>
      {render}
    </Query>
  ),
  prepareWeeklyReview,
  setDueDate,
  setDone,
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
          queryAll: { loading: loadingAll, data: allTasks, error: queryAllError },
          weeklyVisualizationQuery: { loading : weeklyLoading, data: weeklyVisualizationData, error: weeklyError, refetch: weeklyVisualizationRefetch },
          prepareWeeklyReview,
          setDueDate,
          setDone
        }) =>
          <React.Fragment>
            <NavHeader/>
            <Button color='primary' outline size='sm' onClick={() => {
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
            <Button color='primary' outline size='sm' onClick={() => {
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
            </Button>{' '}

            <Container>
              <Row>
                <Col lg={6}>
                  Backlog
                  {loadingAll && <div>loading...</div>}
                  {!loadingAll && console.log('got data', allTasks)}
                  {queryAllError && <div>Tasks: {queryAllError.message}</div>}
                  {(!loadingAll && !queryAllError) && <TaskList listFilter={['Backlog']} setDueDate={setDueDate} setDone={setDone} tasks={allTasks.tasks}/>}
                </Col>
                <Col lg={6}>
                  <Row>
                    Periodic
                    Often
                    {(!loadingAll && !queryAllError) && <TaskList isPeriodic listFilter={['Often']} setDueDate={setDueDate} setDone={setDone} tasks={allTasks.tasks}/>}
                  </Row>
                  <Row>
                    Weekly
                    {(!loadingAll && !queryAllError) && <TaskList noHeader isPeriodic listFilter={['Weekly']} setDueDate={setDueDate} setDone={setDone} tasks={allTasks.tasks}/>}
                  </Row>
                  <Row>
                    Bi-weekly to monthly
                    {(!loadingAll && !queryAllError) && <TaskList noHeader isPeriodic listFilter={['Bi-weekly to monthly']} setDueDate={setDueDate} setDone={setDone} tasks={allTasks.tasks}/>}
                  </Row>
                  <Row>
                    Quarterly to Yearly
                    {(!loadingAll && !queryAllError) && <TaskList noHeader isPeriodic listFilter={['Quarterly to Yearly']} setDueDate={setDueDate} setDone={setDone} tasks={allTasks.tasks}/>}
                  </Row>
                </Col>
              </Row>
              <Row>
                <Col>
                  {weeklyLoading && <div>loading...</div>}
                  {(!weeklyLoading && !weeklyError) &&
                    <div>
                      <FaSync size={25} onClick={() => weeklyVisualizationRefetch()} />
                      <MarkdownRenderer className='weeklyReview' markdown={weeklyVisualizationData.weeklyVisualization} />
                      <style jsx>{`
                      `}</style>
                    </div>
                  }
                  {weeklyError && <div>Weekly review: {weeklyError.message}</div>}
                </Col>
              </Row>
            </Container>
          </React.Fragment>
        }
      </QueryContainer>
    )
  }
}

export default withAlert()(IndexPage)
