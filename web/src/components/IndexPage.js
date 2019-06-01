import React from 'react'
import {
  Button,
  Container,
  Row,
  Col,
  Progress,
  Spinner,
  TabContent,
  TabPane
} from 'reactstrap'
import numeral from 'numeral'

import { Query } from 'react-apollo'
import { adopt } from 'react-adopt'
import { withAlert } from 'react-alert'
import moment from 'moment'
import { 
  FaCalendarDay, 
  FaCheckCircle, 
  FaListUl,
  FaRecycle,
  FaRegClock, 
  FaSync 
} from 'react-icons/fa'
import NavHeader from './NavHeader'
import TaskList from './TaskList'
import GoalList from './GoalList'
import Timer from './Timer'
import ProjectReport from './ProjectReport'

import MarkdownRenderer from 'react-markdown-renderer'

import auth from '../lib/auth0'
// import redirect from '../lib/redirect'
import { navigate } from 'gatsby'

import {
  TaskQuery,
  MonthlyGoalsQuery,
  PrepareWeeklyReview,
  FinishWeeklyReview,
  SetDueDate,
  SetDone,
  SetGoalDone,

  MoveTaskToList,
  WeeklyVisualizationQuery,
  MonthlyVisualizationQuery,
  ActiveTimerQuery,
  StopTimer,
  StartTimer,
  AddTask,
  PrepareMonthlyReview,
  FinishMonthlyReview,
  AddWeeklyGoal,
  AddMonthlyGoal
} from '../lib/graphql'


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
  monthlyVisualizationQuery: ({ render, year, month }) => (
    <Query query={MonthlyVisualizationQuery} ssr={false} variables={{ year, month }}>
      {render}
    </Query>
  ),
  queryTimer: ({ render }) => (
    <Query query={ActiveTimerQuery} ssr={false} variables={{ }}>
      {render}
    </Query>
  ),
  PrepareWeeklyReview,
  FinishWeeklyReview,
  PrepareMonthlyReview,
  FinishMonthlyReview,
  SetDueDate,
  SetDone,
  MoveTaskToList,
  StopTimer,
  StartTimer,
  SetGoalDone,
  AddTask,
  AddWeeklyGoal,
  AddMonthlyGoal,
})

class IndexPage extends React.Component {
  async componentDidMount () {
    console.log('componentDidMount')
    // try {
    //   await auth.instance().silentAuth()
    //   console.log('silentAuth done')
    // } catch (err) {
    //   console.log('error', err)
    //   if (err.error === 'login_required') {
    //     navigate('/login')
    //     // redirect({}, '/login')
    //     return null
    //   }
    // }
    if (!auth.instance().isAuthenticated()) {
      console.log('not authenticated')
      navigate('/app/login')
      return null
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
    let nowGraceMonth = moment().subtract(5,'days')
    let monthNext = nowGraceMonth.month()+1
    let elapsedWeek = moment().diff(moment().startOf('week').hours(12), 'hours', true)
    let remainingWeek = 168 - elapsedWeek
    let elapsedWeekPct = elapsedWeek/168.0*100.0
    console.log('IndexPage props are', this.props)
    return (
      <QueryContainer year={now.year()} week={now.isoWeek()} month={now.month()+1}>
        {({
          queryAll: { loading: loadingAll, data: allTasks, error: queryAllError, refetch: allRefetch },
          queryAllGoals: { loading: loadingAllGoals, data: allGoals, error: queryAllGoalsError, refetch: allGoalsRefetch },
          weeklyVisualizationQuery: { loading : weeklyLoading, data: weeklyVisualizationData, error: weeklyError, refetch: weeklyVisualizationRefetch },
          monthlyVisualizationQuery: { loading : monthlyLoading, data: monthlyVisualizationData, error: monthlyError, refetch: monthlyVisualizationRefetch },
          queryTimer: { loading: loadingTimer, data: timerData, error: timerError, refetch: timerRefetch },
          PrepareWeeklyReview,
          FinishWeeklyReview,
          PrepareMonthlyReview,
          FinishMonthlyReview,
          SetDueDate,
          SetDone,
          MoveTaskToList,
          StopTimer,
          StartTimer,
          SetGoalDone,
          AddTask,
          AddWeeklyGoal,
          AddMonthlyGoal,
        }) =>
          <React.Fragment>
            <NavHeader switchTab={this.switchTab} activeTab={this.state.activeTab} />


            <Container>
              <TabContent activeTab={this.state.activeTab}>
                <TabPane tabId="board">
                  <Progress className="weeklyProgress" color={'info'} value={elapsedWeekPct}>
                    {`${numeral(remainingWeek).format('0')} hours remaining until weekly review (Sun 12pm)`}
                  </Progress>
                  <FaSync size={25} onClick={() => {
                    allRefetch()
                    allGoalsRefetch()
                    timerRefetch()
                  }} />

                  <Row>
                    <Col lg={6}>
                      {loadingTimer && <Spinner color="primary" />}
                      {!loadingTimer && console.log('got data', timerData)}
                      {timerError && <div>Timer: {timerError.message}</div>}
                      {!loadingTimer && !timerError && <Timer stopTimer={StopTimer} timerRefetch={timerRefetch} activeTimer={timerData.activeTimer} />}
                    </Col>
                    <Col lg={6}>
                      <TaskList
                        loading={loadingAll} error={queryAllError} data={allTasks}
                        listTitle={<div><FaCalendarDay size={25}/> {'Today'}</div>}
                        listFilter={['Today']}
                        setDueDate={SetDueDate} setDone={SetDone}
                        moveTaskToList={MoveTaskToList} startTimer={StartTimer}
                        timerRefetch={timerRefetch}
                        addTask={AddTask} board={'Kanban daily/weekly'} list={'Today'}
                      />

                    </Col>
                  </Row>
                  <Row>
                    <Col lg={6}>
                      <GoalList 
                        loading={loadingAllGoals} error={queryAllGoalsError}
                        startTimer={StartTimer} timerRefetch={timerRefetch}
                        setGoalDone={SetGoalDone}
                        addWeeklyGoal={AddWeeklyGoal}
                        addMonthlyGoal={AddMonthlyGoal}
                        goals={allGoals ? allGoals.monthlyGoals : null}
                      />
                    </Col>

                  </Row>
                  <Row>
                    <Col lg={6}>
                      <TaskList
                        loading={loadingAll} error={queryAllError} data={allTasks}
                        listTitle={<div><FaListUl size={25}/> {'Backlog'}</div>}
                        listFilter={['Backlog (Personal)']}
                        setDueDate={SetDueDate} setDone={SetDone}
                        moveTaskToList={MoveTaskToList} startTimer={StartTimer}
                        timerRefetch={timerRefetch}
                        addTask={AddTask} board={'Backlog (Personal)'} list={'Backlog'}
                      />
                    </Col>
                    <Col lg={6}>
                      <TaskList
                        loading={loadingAll} error={queryAllError} data={allTasks}
                        listTitle={<div><FaRegClock size={25}/> {'Waiting on...'}</div>}
                        listFilter={['Waiting on']}
                        setDueDate={SetDueDate} setDone={SetDone}
                        moveTaskToList={MoveTaskToList} startTimer={StartTimer}
                        timerRefetch={timerRefetch}
                        addTask={AddTask} board={'Kanban daily/weekly'} list={'Waiting on'}
                      />
                    </Col>
                  </Row>
                  <Row>
                    <Col lg={6}>
                      <TaskList
                        loading={loadingAll} error={queryAllError} data={allTasks}
                        listTitle={<div><FaCheckCircle size={25}/> {'Done this week'}</div>}
                        listFilter={['Done this week']}
                        setDueDate={SetDueDate} setDone={SetDone}
                        moveTaskToList={MoveTaskToList} startTimer={StartTimer}
                        timerRefetch={timerRefetch}
                      />
                    </Col>
                  </Row>
                </TabPane>
                <TabPane tabId="periodicBoard">
                  <FaSync size={25} onClick={() => {
                    allRefetch()
                  }} />
                  <Row>
                    <Col lg={6}>
                      <Row>
                        <div className="listTitle"><FaRecycle size={25}/> {'Periodic'}</div>
                        <TaskList
                          loading={loadingAll} error={queryAllError} data={allTasks}
                          listSubGroupTitle={'Often'} listFilter={['Often']}
                          setDueDate={SetDueDate} setDone={SetDone}
                          moveTaskToList={MoveTaskToList} startTimer={StartTimer}
                          timerRefetch={timerRefetch}
                        />
                      </Row>
                      <Row>
                        <TaskList
                          loading={loadingAll} error={queryAllError} data={allTasks}
                          listSubGroupTitle={'Weekly'} listFilter={['Weekly']}
                          setDueDate={SetDueDate} setDone={SetDone}
                          moveTaskToList={MoveTaskToList} startTimer={StartTimer}
                          timerRefetch={timerRefetch}
                        />
                      </Row>
                      <Row>
                        <TaskList
                          loading={loadingAll} error={queryAllError} data={allTasks}
                          listSubGroupTitle={'Bi-weekly to monthly'} listFilter={['Bi-weekly to monthly']}
                          setDueDate={SetDueDate} setDone={SetDone}
                          moveTaskToList={MoveTaskToList} startTimer={StartTimer}
                          timerRefetch={timerRefetch}
                        />
                      </Row>
                    </Col>
                    <Col lg={6}>
                      <Row>
                        <TaskList
                          loading={loadingAll} error={queryAllError} data={allTasks}
                          listSubGroupTitle={'Quarterly to Yearly'} listFilter={['Quarterly to Yearly']}
                          setDueDate={SetDueDate} setDone={SetDone}
                          moveTaskToList={MoveTaskToList} startTimer={StartTimer}
                          timerRefetch={timerRefetch}
                        />
                      </Row>
                    </Col>
                  </Row>
                </TabPane>
                <TabPane tabId="weeklyReview">
                  <Button color='primary' size='sm' onClick={() => {
                    PrepareWeeklyReview
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
                      PrepareWeeklyReview
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
                    FinishWeeklyReview
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
                    <Col>
                      {weeklyLoading && <Spinner color="primary" />}
                      <div>
                        {!weeklyLoading && <FaSync size={25} onClick={() => weeklyVisualizationRefetch()} />}
                        {(!weeklyLoading && !weeklyError) &&
                            <MarkdownRenderer className='weeklyReview' markdown={weeklyVisualizationData.weeklyVisualization} />
                        }
                      </div>
                      {weeklyError && <div>Weekly review: {weeklyError.message}</div>}
                    </Col>
                  </Row>
                </TabPane>
                <TabPane tabId="monthlyReview">
                  <Button color='primary' size='sm' onClick={() => {
                    PrepareMonthlyReview
                      .mutation({
                        variables: {
                          year: nowGraceMonth.year(),
                          month: nowGraceMonth.month()+1,
                        }
                      })
                      .then(({ data }) => alert.show(data.prepareMonthlyReview.message))
                  }}>
                      Prepare monthly review for {nowGraceMonth.year()}-{nowGraceMonth.month()+1}
                  </Button>{' '}
                  {nowGraceMonth.month()+1 !== monthNext &&
                    <Button color='primary' size='sm' onClick={() => {
                      PrepareMonthlyReview
                        .mutation({
                          variables: {
                            year: nowGraceMonth.year(),
                            month: monthNext,
                          }
                        })
                        .then(({ data }) => alert.show(data.prepareMonthlyReview.message))
                    }}>
                        Prepare monthly review for {now.year()}-{monthNext}
                    </Button>
                  }{' '}
                  {nowGraceMonth.month()+1 !== monthNext &&
                    <Button color='primary' size='sm' onClick={() => {
                      FinishMonthlyReview
                        .mutation({
                          variables: {
                            year: nowGraceMonth.year(),
                            month: nowGraceMonth.month()+1,
                          }
                        })
                        .then(({ data }) => alert.show(data.finishMonthlyReview.message))
                    }}>
                        Finish monthly review for {nowGraceMonth.year()}-{nowGraceMonth.month()+1}
                    </Button>
                  }{' '}
                  <Button color='primary' size='sm' onClick={() => {
                    FinishMonthlyReview
                      .mutation({
                        variables: {
                          year: nowGraceMonth.year(),
                          month: monthNext,
                        }
                      })
                      .then(({ data }) => alert.show(data.finishMonthlyReview.message))
                  }}>
                      Finish monthly review for {nowGraceMonth.year()}-{monthNext}
                  </Button>{' '}
                  <Row>
                    <Col>
                      {monthlyLoading && <Spinner color="primary" />}
                      <div>
                        {!monthlyLoading && <FaSync size={25} onClick={() => monthlyVisualizationRefetch()} />}
                        {(!monthlyLoading && !monthlyError) &&
                          <MarkdownRenderer className='monthlyReview' markdown={monthlyVisualizationData.monthlyVisualization} />
                        }
                      </div>
                      {monthlyError && <div>Monthly review: {monthlyError.message}</div>}
                    </Col>
                  </Row>
                </TabPane>
                <TabPane tabId="timeReport">
                  <ProjectReport />
                </TabPane>
              </TabContent>
            </Container>
            <style jsx global>{`
              .weeklyProgress {
                margin-bottom: 1em;
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

export default IndexPage
