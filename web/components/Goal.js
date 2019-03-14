import React from 'react'
import {
  Row,
  Col,
} from 'reactstrap'
import moment from 'moment'

class Goal extends React.Component {
  constructor (props) {
    super(props)
  }
  render () {
    let { goal } = this.props
    let now = moment()
    let thisWeek = now.isoWeek()
    return (
      <React.Fragment key={goal.id}>
        <Row>
          <Col>Monthly: {goal.title}</Col>
        </Row>
        {goal.weeklyGoals
          .filter(g => g.week === thisWeek)
          .sort((a,b) => b.week - a.week)
          .map((g)=>
            <Row key={g.title+g.week}>
              <Col>{g.week}: {g.title}</Col>
            </Row>
          )
        }
      </React.Fragment>
    )
  }
}

export default Goal
