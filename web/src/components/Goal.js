import React from 'react'
import {
  Button,
  Collapse,
  Row,
  Col,
  Form,
  Input
} from 'reactstrap'
import { 
  FaStarHalfAlt,
  FaStar,
  FaStop
} from 'react-icons/fa'

import moment from 'moment'

class Goal extends React.Component {
  constructor (props) {
    super(props)
    let now = moment()
    let thisWeek = now.isoWeek()
    this.state = {
      showControls: false,
      showAddControls: false,
      title: '',
      week: thisWeek,
    }
    this.toggle = this.toggle.bind(this)
    this.toggleAdd = this.toggleAdd.bind(this)
    this.handleWeekChange = this.handleWeekChange.bind(this)
    this.handleTitleChange = this.handleTitleChange.bind(this)
  }
  toggle() {
    this.setState(state => ({ showControls: !state.showControls }))
  }
  toggleAdd() {
    this.setState(state => ({ showAddControls: !state.showAddControls }))
  }
  handleTitleChange(e) {
    this.setState({ title: e.target.value })
  }
  handleWeekChange(e) {
    this.setState({ week: e.target.value })
  }
  render () {
    let { goal, startTimer, timerRefetch, setGoalDone, addWeeklyGoal } = this.props
    let now = moment()
    let thisWeek = now.isoWeek()
    return (
      <React.Fragment key={goal.id}>
        <Row className='monthlyGoal'>
          <Col onClick={this.toggleAdd}>Monthly: {goal.title}</Col>
          <style jsx global>{`
              .monthlyGoal {
                background: #999;
              }
              .goal-done {
                text-decoration: line-through;
              }
              #newGoalWeek {
                width: 3em;
                display: inline-block;
              }
              #newGoalTitle {
                width: 20em;
                display: inline-block;
              }
            `}</style>
        </Row>
        <Collapse isOpen={this.state.showAddControls}>
          <Form>
            {`Week:`}<Input type="text"
              value={this.state.week}
              size={'5'}
              id="newGoalWeek"
              placeholder="week"
              onChange={this.handleWeekChange}/>
            <Input type="text"
              value={this.state.title}
              size={'50'}
              id="newGoalTitle"
              placeholder="goal description"
              onChange={this.handleTitleChange}/>
            <Button key={'add'+goal.title} onClick={() => {
              addWeeklyGoal.mutation({
                variables: {
                  taskID: goal.idCard,
                  title: this.state.title,
                  week: this.state.week,
                }
              })
                .then(() => {
                  this.setState({ showAddControls: false })
                })
            }} size="sm" color="primary">
              Add goal</Button>
          </Form>
        </Collapse>
        {goal.weeklyGoals
          .filter(g => (g.week === thisWeek || g.week === thisWeek+1))
          .sort((a,b) => b.week - a.week)
          .map((g)=> {
            let doneClass = ''
            if (g.done) {
              doneClass='goal-done'
            }
            return (
              <Row key={g.idCard+g.idCheckitem}>
                <Col>
                  <div className={`goal ${doneClass}`} onClick={this.toggle}>
                  {g.done && g.status == '(not done)' && <FaStop size={25} color={'red'}/>}
                  {g.done && g.status == '(done)' && <FaStar size={25}/>}
                  {g.done && g.status == '(partial)' && <FaStarHalfAlt size={25}/>}{' '}{g.week}: {g.title} {g.status}</div>
                  <Collapse isOpen={this.state.showControls}>
                    <Button outline color='primary' size='sm' onClick={()=>{
                      startTimer.mutation({
                        variables: {
                          taskID: g.idCard,
                          checkitemID: g.idCheckitem
                        }
                      })
                        .then(() => timerRefetch())
                    }}>Start</Button>
                    <Button outline color='primary' size='sm' onClick={()=>{
                      setGoalDone.mutation({
                        variables: {
                          taskId: g.idCard,
                          checkitemID: g.idCheckitem,
                          done: true,
                          status: '(done)'
                        }
                      })
                    }}>√</Button>
                    <Button outline color='primary' size='sm' onClick={()=>{
                      setGoalDone.mutation({
                        variables: {
                          taskId: g.idCard,
                          checkitemID: g.idCheckitem,
                          done: true,
                          status: '(partial)',
                        }
                      })
                    }}>½√</Button>
                    <Button outline color='primary' size='sm' onClick={()=>{
                      setGoalDone.mutation({
                        variables: {
                          taskId: g.idCard,
                          checkitemID: g.idCheckitem,
                          done: true,
                          status: '(not done)'
                        }
                      })
                    }}>X√</Button>
                    <Button outline color='primary' size='sm' onClick={()=>{
                      setGoalDone.mutation({
                        variables: {
                          taskId: g.idCard,
                          checkitemID: g.idCheckitem,
                          done: false,
                          status: '',
                        }
                      })
                    }}>X</Button>
                  </Collapse>
                </Col>
              </Row>
            )
          })
        }
      </React.Fragment>
    )
  }
}

export default Goal
