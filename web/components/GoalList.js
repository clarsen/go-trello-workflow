import React from 'react'
import Goal from './Goal'
import { Container } from 'reactstrap'

class GoalList extends React.Component {
  constructor (props) {
    super(props)
  }
  render () {
    let { goals, startTimer, timerRefetch } = this.props
    return (
      <Container>
        {
          goals
            .map((g) => <Goal key={g.title} startTimer={startTimer} timerRefetch={timerRefetch} goal={g}/>)
        }
      </Container>
    )
  }
}

export default GoalList
