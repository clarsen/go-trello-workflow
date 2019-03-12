import React from 'react'
import {
  Button,
  Collapse,
  Row,
  Col,
  Progress
} from 'reactstrap'
import Moment from 'react-moment'
import moment from 'moment'
import { MdCheckBoxOutlineBlank } from 'react-icons/md'

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
    return (
      <React.Fragment>
        <Row>
          {//<Col xs={2} lg={2}>Created</Col>
          }
          { // <Col xs={1} lg={1}>List</Col>
          }
          <Col xs={9} lg={9}>Title</Col>
          <Col xs={3} lg={3}>Due</Col>
        </Row>
      </React.Fragment>
    )
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
          {//<Col xs={2} lg={2} key={'1'+task.id}><Moment unix fromNow withTitle titleFormat={'LL'}>{task.createdDate}</Moment></Col>
          }
          <Col xs={9} lg={9} key={'3'+task.id}>
            {!task.due && <MdCheckBoxOutlineBlank
              size={20}
              onClick={() => {
                setDone.mutation({
                  variables: {
                    taskId: task.id,
                    done: true,
                  }
                })
              }}
            />}
            <a target="_blank" rel="noopener noreferrer" href={task.url}>
              {task.title}
            </a></Col>
          <Col xs={2} lg={3} key={'2'+task.id}>
            { task.due &&
              <div>
                {value>0
                  ?
                  <Progress color={color} value={value} onClick={this.toggle}>
                    <Moment unix fromNow withTitle titleFormat={'LL'}>{task.due}</Moment>
                  </Progress>
                  :
                  <Moment unix fromNow withTitle titleFormat={'LL'} onClick={this.toggle}>{task.due}</Moment>
                }
                <Collapse isOpen={this.state.showDueDateControls}>
                  <Button outline color='primary' size='sm' onClick={()=>{
                    setDone.mutation({
                      variables: {
                        taskId: task.id,
                        done: true,
                        nextDue: moment.unix(task.due).add(1, 'month').unix(),
                      }
                    })
                  }}>√+1m</Button>{' '}
                  <Button outline color='primary' size='sm' onClick={()=>{
                    setDueDate.mutation({
                      variables: {
                        taskId: task.id,
                        due: task.due + 7*86400
                      }
                    })
                  }}>+=1w</Button>{' '}
                  <Button outline color='primary' size='sm' onClick={()=>{
                    setDueDate.mutation({
                      variables: {
                        taskId: task.id,
                        due: moment().add(7, 'days').unix(),
                      }
                    })
                  }}>+1w</Button>{' '}
                  <Button outline color='primary' size='sm' onClick={()=>{
                    setDueDate.mutation({
                      variables: {
                        taskId: task.id,
                        due: moment().add(1, 'months').unix(),
                      }
                    })
                  }}>+1m</Button>{' '}
                </Collapse>
              </div>
            }
          </Col>
          { //<Col xs={1} lg={1} key={'list'+task.id}>{task.list}</Col>
          }
        </Row>
      </React.Fragment>

    )
  }
}

export default Task
