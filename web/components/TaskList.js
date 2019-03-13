import React from 'react'
import Task from './Task'
import { Container } from 'reactstrap'

class TaskList extends React.Component {
  constructor (props) {
    super(props)
  }
  render () {
    let { listFilter, noHeader, isPeriodic, tasks, setDueDate, setDone } = this.props
    return (
      <Container>
        {!noHeader && Task.header()}
        {
          tasks
            .filter((t) => !isPeriodic || isPeriodic && t.period)
            .filter((t) => !listFilter || listFilter.count == 0 || listFilter.includes(t.list.list))
            .sort((a,b) => b.createdDate - a.createdDate)
            .sort((a,b) => {
              if (a.due && b.due) { // within tasks with due dates, earlier ones first
                return a.due - b.due
              } else if (a.due) { // tasks without due dates sorted before those with due dates
                return 1
              } else if (b.due ) { // tasks without due dates sorted before those with due dates
                return -1
              } else { // within tasks with no due dates, defer to created date
                return 0
              }
            })
            .map((t) => <Task key={t.id} setDueDate={setDueDate} setDone={setDone} task={t}/>)
        }
      </Container>
    )
  }
}

export default TaskList
