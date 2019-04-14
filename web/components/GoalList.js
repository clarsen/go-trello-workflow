import React from 'react'
import Goal from './Goal'
import {
  Container
} from 'reactstrap'

class GoalList extends React.Component {
  constructor (props) {
    super(props)
  }
  render () {
    let { goals, startTimer, timerRefetch, setGoalDone, addWeeklyGoal } = this.props
    return (
      <Container>
        {
          goals
            .map((g) => <Goal key={g.title} startTimer={startTimer} setGoalDone={setGoalDone} timerRefetch={timerRefetch} addWeeklyGoal={addWeeklyGoal} goal={g}/>)
        }
      </Container>
    )
  }
}

export default GoalList
