import React from 'react'
import { Row, Col } from 'reactstrap'

class Task extends React.Component {
  constructor (props) {
    super(props)
  }
  render () {
    let { task } = this.props
    return (
      <React.Fragment key={task.id}>
        <Row key={'row0'+task.id}>
          <Col xs='auto' key={'1'+task.id}>{task.id}</Col>
          <Col xs='auto' key={'2'+task.id}>{task.title}</Col>
        </Row>
      </React.Fragment>
    )
  }
}

export default Task
