import React from 'react'
import {
  Button,
  Container,
  Row,
  Col,
  Spinner,
  TabContent,
  TabPane
} from 'reactstrap'

import { Query, Mutation } from 'react-apollo'
import { adopt } from 'react-adopt'
import { withAlert } from 'react-alert'
import moment from 'moment'
import { FaSync } from 'react-icons/fa'
import NavHeader from '../components/NavHeader'
import TaskList from '../components/TaskList'
import GoalList from '../components/GoalList'
import Timer from '../components/Timer'
import MarkdownRenderer from 'react-markdown-renderer'

import auth from '../lib/auth0'
import redirect from '../lib/redirect'

import {
  TaskQuery,
  MonthlyGoalsQuery,
  PrepareWeeklyReviewQuery,
  FinishWeeklyReviewQuery,
  SetDueDateQuery,
  SetDoneQuery,
  MoveTaskToListQuery,
  WeeklyVisualizationQuery,
  ActiveTimerQuery,
  StopTimerQuery,
  StartTimerQuery,
  SetGoalDoneQuery
} from '../lib/graphql'


const prepareWeeklyReview = ({ render }) => (
  <Mutation
    mutation={PrepareWeeklyReviewQuery}
  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)


const finishWeeklyReview = ({ render }) => (
  <Mutation
    mutation={FinishWeeklyReviewQuery}
  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)

const setDueDate = ({ render }) => (
  <Mutation
    mutation={SetDueDateQuery}
  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)


const setDone = ({ render }) => (
  <Mutation
    mutation={SetDoneQuery}
  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)

const setGoalDone = ({ render }) => (
  <Mutation
    mutation={SetGoalDoneQuery}
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


const moveTaskToList = ({ render }) => (
  <Mutation
    mutation={MoveTaskToListQuery}
  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)


const stopTimer = ({ render }) => (
  <Mutation
    mutation={StopTimerQuery}
  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)

const startTimer = ({ render }) => (
  <Mutation
    mutation={StartTimerQuery}
  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)



/* eslint-disable react/display-name */
/* eslint-disable react/prop-types */

const QueryContainer = adopt({
  queryAll: ({ render }) => (
    <Query query={TaskQuery} ssr={false} variables={{ }}>
      {render}
    </Query>
  ),
  queryAllGoals: ({ render }) => (
    <Query query={MonthlyGoalsQuery} ssr={false} variables={{ }}>
      {render}
    </Query>
  ),
  weeklyVisualizationQuery: ({ render, year, week }) => (
    <Query query={WeeklyVisualizationQuery} ssr={false} variables={{ year, week }}>
      {render}
    </Query>
  ),
  queryTimer: ({ render }) => (
    <Query query={ActiveTimerQuery} ssr={false} variables={{ }}>
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
  setGoalDone,
})

class IndexPage extends React.Component {
  componentDidMount () {
    console.log('componentDidMount')
    if (!auth().isAuthenticated()) {
      console.log('not authenticated')
      redirect({}, '/login')
    }
  }
  constructor (props) {
    super(props)
    this.state = {
      data: null,
      activeTab: 'board',
    }
    this.switchTab = this.switchTab.bind(this)
  }
  switchTab (tab) {
    if (this.state.activeTab !== tab) {
      this.setState({
        activeTab: tab
      })
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
          setGoalDone,
        }) =>
          <React.Fragment>
            <NavHeader switchTab={this.switchTab} activeTab={this.state.activeTab} />


            <Container>
              <TabContent activeTab={this.state.activeTab}>
                <TabPane tabId="board">
                  <FaSync size={25} onClick={() => {
                    allRefetch()
                    allGoalsRefetch()
                    timerRefetch()
                  }} />
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
                      .then(({ data }) => {
                        alert.show(data.finishWeeklyReview.message)
                        allRefetch()
                        allGoalsRefetch()
                      })
                  }}>
                      Finish weekly review for {now.year()}-{now.isoWeek()}
                  </Button>{' '}
                  <Row>
                    <Col lg={6}>
                      {loadingTimer && <Spinner color="primary" />}
                      {!loadingTimer && console.log('got data', timerData)}
                      {timerError && <div>Timer: {timerError.message}</div>}
                      {!loadingTimer && !timerError && <Timer stopTimer={stopTimer} timerRefetch={timerRefetch} activeTimer={timerData.activeTimer} />}
                    </Col>
                    <Col lg={6}>
                      <div className="listTitle">Today</div>
                      {loadingAll && <Spinner color="primary" />}
                      {!loadingAll && console.log('got data', allTasks)}
                      {queryAllError && <div>Tasks: {queryAllError.message}</div>}
                      {(!loadingAll && !queryAllError) && <TaskList listFilter={['Today']} setDueDate={setDueDate} setDone={setDone} moveTaskToList={moveTaskToList} startTimer={startTimer} timerRefetch={timerRefetch} tasks={allTasks.tasks}/>}
                    </Col>
                  </Row>
                  <Row>
                    <Col lg={6}>
                      <div className="listTitle">Goals</div>
                      {loadingAllGoals && <Spinner color="primary" />}
                      {!loadingAllGoals && console.log('got data', allGoals)}
                      {queryAllGoalsError && <div>Goals: {queryAllGoalsError.message}</div>}
                      {!loadingAllGoals && !queryAllGoalsError && <GoalList startTimer={startTimer} timerRefetch={timerRefetch} setGoalDone={setGoalDone} goals={allGoals.monthlyGoals}/>}
                    </Col>
                    <Col lg={6}>
                      <div className="listTitle">Waiting on...</div>
                      {loadingAll && <Spinner color="primary" />}
                      {(!loadingAll && !queryAllError) && <TaskList listFilter={['Waiting on']} setDueDate={setDueDate} setDone={setDone} moveTaskToList={moveTaskToList}  startTimer={startTimer} timerRefetch={timerRefetch} tasks={allTasks.tasks}/>}
                    </Col>
                  </Row>
                  <Row>
                    <Col lg={6}>
                      <div className="listTitle">Done this week</div>
                      {loadingAll && <Spinner color="primary" />}
                      {(!loadingAll && !queryAllError) && <TaskList listFilter={['Done this week']} setDueDate={setDueDate} setDone={setDone} moveTaskToList={moveTaskToList}  startTimer={startTimer} timerRefetch={timerRefetch} tasks={allTasks.tasks}/>}
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
                </TabPane>
                <TabPane tabId="weeklyReview">
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
                </TabPane>
              </TabContent>
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
