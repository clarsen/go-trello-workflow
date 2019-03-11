import React from 'react'
import { Row, Col } from 'reactstrap'
import Moment from 'react-moment'

class Task extends React.Component {
  constructor (props) {
    super(props)
  }
  render () {
    let { task } = this.props
    return (
      <React.Fragment key={task.id}>
        <Row key={'row0'+task.id}>
          <Col xs='auto' key={'2'+task.id}><Moment unix fromNow withTitle titleFormat={'LL'}>{task.createdDate}</Moment></Col>
          <Col xs='auto' key={'3'+task.id}><a target="_blank" rel="noopener noreferrer" href={task.url}>{task.title}</a></Col>
        </Row>
      </React.Fragment>
    )
  }
}

export default Task
