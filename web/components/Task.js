import React from 'react'
import {
  Button,
  Collapse,
  Row,
  Col,
  Progress
} from 'reactstrap'
import moment from 'moment'

class Task extends React.Component {
  constructor (props) {
    super(props)
    this.state = { showDueDateControls: false }
    this.toggle = this.toggle.bind(this)
  }
  toggle() {
    this.setState(state => ({ showDueDateControls: !state.showDueDateControls }))
  }
  static header() {
    return null
  }
  render () {
    let { task, setDueDate, setDone } = this.props
    let color = ''
    let value = 0
    if (task.due) {
      let delta_days = moment().diff(moment.unix(task.due))/(86400*1000)
      console.log('task', task.title, 'task.due', moment.unix(task.due), 'delta_days from now', delta_days)
      if (delta_days < -3) {
        color = 'info'
        value = (100*(14 + delta_days)/14).toFixed(0)
        if (value < 0) {
          value = 0
        }
      } else if (delta_days >= -3 && delta_days < 0) {
        color = 'warning'
        value = (100*(14 + delta_days)/14).toFixed(0)
      } else if (delta_days >= 0){
        color = 'danger'
        value= ((delta_days/7)*100).toFixed(0)
        if (value > 100) {
          value = 100
        }
      }
    }
    return (
      <React.Fragment key={task.id}>
        <Row key={'row0'+task.id}>
          <Col xs={12} lg={12} key={'2'+task.id}>
            <div>
              <div className='task' onClick={this.toggle}>
                {task.title}{' '}
                { (task.due && value>0) &&
                      <Progress color={color} value={value} onClick={this.toggle}>
                        {moment.unix(task.due).fromNow()}
                      </Progress>
                }
                { (task.due && value==0) && moment.unix(task.due).fromNow() }
              </div>
              <style jsx global>{`
                .progress {
                  float: right;
                  background-color: #888;
                }
              `}</style>

              <Collapse isOpen={this.state.showDueDateControls}>
                {!task.period &&
                  <Button outline color='primary' size='sm' onClick={()=>{
                    setDone.mutation({
                      variables: {
                        taskId: task.id,
                        done: true,
                      }
                    })
                  }}>√</Button>
                }{' '}
                {task.period &&
                  <Button outline color='primary' size='sm' onClick={()=>{
                    setDone.mutation({
                      variables: {
                        taskId: task.id,
                        done: true,
                        nextDue: moment.unix(task.due).add(1, 'month').unix(),
                      }
                    })
                  }}>√+1m</Button>
                }{' '}
                {task.due &&
                  <Button outline color='primary' size='sm' onClick={()=>{
                    setDueDate.mutation({
                      variables: {
                        taskId: task.id,
                        due: task.due + 7*86400
                      }
                    })
                  }}>+=1w</Button>
                }{' '}
                {task.due &&
                  <Button outline color='primary' size='sm' onClick={()=>{
                    setDueDate.mutation({
                      variables: {
                        taskId: task.id,
                        due: moment().add(7, 'days').unix(),
                      }
                    })
                  }}>+1w</Button>
                }{' '}
                {task.due &&
                  <Button outline color='primary' size='sm' onClick={()=>{
                    setDueDate.mutation({
                      variables: {
                        taskId: task.id,
                        due: moment().add(1, 'months').unix(),
                      }
                    })
                  }}>+1m</Button>
                }{' '}
              </Collapse>
            </div>
          </Col>
        </Row>
      </React.Fragment>

    )
  }
}

export default Task
