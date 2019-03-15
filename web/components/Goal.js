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
    let { goal, startTimer, timerRefetch } = this.props
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
            `}</style>
        </Row>
        {goal.weeklyGoals
          .filter(g => g.week === thisWeek)
          .sort((a,b) => b.week - a.week)
          .map((g)=>
            <Row key={g.idCard+g.idCheckitem}>
              <Col>
                <div className='goal' onClick={this.toggle}>{g.week}: {g.title}</div>
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
                </Collapse>
              </Col>
            </Row>
          )
        }
      </React.Fragment>
    )
  }
}

export default Goal
