import React from 'react'
import TaskReviewItem from './TaskReviewItem'
import {
  Button,
  Collapse,
  Container,
  Form,
  Input,
  Spinner,
} from 'reactstrap'
import moment from 'moment'

class TaskReviewList extends React.Component {
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
    let { listTitle, loading, error, data,
      boardFilter, listFilter, 
      setDone
    } = this.props
    // console.log('for list', list)
    return (
      <React.Fragment>
        {listTitle && <div className="listTitle" onClick={this.toggle}>{listTitle}</div>}
        {loading && <Spinner color="primary" />}
        {error && <div>Tasks: {error.message}</div>}
        {(!loading && !error) &&
          <Container>
            {
              data.tasks
                .filter((t) => (!listFilter || listFilter.count == 0 || listFilter.includes(t.list.list)) && (!boardFilter || boardFilter.count == 0 || boardFilter.includes(t.list.board))) 
                .sort((a,b) => a.createdDate - b.createdDate)
                .sort((a,b) => a.dateLastActivity - b.dateLastActivity)
                .map((t) => <TaskReviewItem key={t.id} setDone={setDone} task={t}/>)
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

export default TaskReviewList
