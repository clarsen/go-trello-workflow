import React from 'react'
import Task from './Task'
import {
  Button,
  Collapse,
  Container,
  Form,
  Input,
  Spinner,
} from 'reactstrap'

class TaskList extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      showControls: false,
      title: ''
    }
    this.toggle = this.toggle.bind(this)
    this.handleChange = this.handleChange.bind(this)
  }
  toggle() {
    this.setState(state => ({ showControls: !state.showControls }))
  }
  handleChange(e) {
    this.setState({ title: e.target.value })
  }

  render () {
    let { listTitle, listSubGroupTitle, loading, error, data,
      board, list,
      listFilter, isPeriodic, setDueDate, setDone, moveTaskToList, startTimer, timerRefetch,
      addTask,
    } = this.props
    return (
      <React.Fragment>
        {listTitle && <div className="listTitle" onClick={this.toggle}>{listTitle}</div>}
        {listSubGroupTitle && <div className="listSubGroupTitle" onClick={this.toggle}>{listSubGroupTitle}</div>}
        {addTask &&
          <Collapse isOpen={this.state.showControls}>
            <br/>
            <Form>
              <Input type="text"
                value={this.state.title}
                size={'50'}
                id="newTaskTitle"
                placeholder="task description"
                onChange={this.handleChange}/>
              <Button key={'add'+listTitle+listSubGroupTitle} onClick={() => {
                addTask.mutation({
                  variables: {
                    title: this.state.title,
                    board: board,
                    list: list,
                  }
                })
                  .then(() => {
                    this.setState({ showControls: false })
                  })
              }} size="sm" color="primary">
                Add task</Button>
            </Form>
          </Collapse>
        }
        {loading && <Spinner color="primary" />}
        {!loading && console.log('got data', data)}
        {error && <div>Tasks: {error.message}</div>}
        {(!loading && !error) &&
          <Container>
            {
              data.tasks
                .filter((t) => !isPeriodic || isPeriodic && t.period)
                .filter((t) => !listFilter || listFilter.count == 0 || listFilter.includes(t.list.board) || listFilter.includes(t.list.list))
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
                .map((t) => <Task key={t.id} setDueDate={setDueDate} setDone={setDone} moveTaskToList={moveTaskToList} startTimer={startTimer} timerRefetch={timerRefetch} task={t}/>)
            }
          </Container>
        }
        <style global jsx>{`
          #newTaskTitle {
            width: 50%;
          }
          .listSubGroupTitle {
            background: #999;
            width: 100%;
          }
          .listTitle {
            background: #bbb;
            width: 100%;
            color: #fff;
          }
        `}</style>
      </React.Fragment>
    )
  }
}

export default TaskList
