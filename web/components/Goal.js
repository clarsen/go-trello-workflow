import React from 'react'
import {
  Button,
  Collapse,
  Row,
  Col,
} from 'reactstrap'
import moment from 'moment'

class Goal extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      showControls: false,
    }
    this.toggle = this.toggle.bind(this)
  }
  toggle() {
    this.setState(state => ({ showControls: !state.showControls }))
  }
  render () {
    let { goal, startTimer, timerRefetch, setGoalDone } = this.props
    let now = moment()
    let thisWeek = now.isoWeek()
    return (
      <React.Fragment key={goal.id}>
        <Row className='monthlyGoal'>
          <Col>Monthly: {goal.title}</Col>
          <style jsx global>{`
              .monthlyGoal {
                background: #999;
              }
              .goal-done {
                text-decoration: line-through;
              }
            `}</style>
        </Row>
        {goal.weeklyGoals
          .filter(g => g.week === thisWeek)
          .sort((a,b) => b.week - a.week)
          .map((g)=> {
            let doneClass = ''
            if (g.done) {
              doneClass='goal-done'
            }
            return (
              <Row key={g.idCard+g.idCheckitem}>
                <Col>
                  <div className={`goal ${doneClass}`} onClick={this.toggle}>{g.week}: {g.title} {g.status}</div>
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
