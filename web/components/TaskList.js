import React from 'react'
import Task from './Task'
import { Container } from 'reactstrap'

class TaskList extends React.Component {
  constructor (props) {
    super(props)
  }
  render () {
    let { tasks } = this.props
    return (
      <Container>
        {
          tasks.map((t) => <Task key={t.id} task={t}/>)
        }
      </Container>
    )
  }
}

export default TaskList
