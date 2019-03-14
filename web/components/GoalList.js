import React from 'react'
import Goal from './Goal'
import { Container } from 'reactstrap'

class GoalList extends React.Component {
  constructor (props) {
    super(props)
  }
  render () {
    let { goals } = this.props
    return (
      <Container>
        {
          goals
            .map((g) => <Goal key={g.id} goal={g}/>)
        }
      </Container>
    )
  }
}

export default GoalList
