import React from 'react'
import Goal from './Goal'
import {
  Button,
  Collapse,
  Container,
  Form,
  Input,
  Spinner
} from 'reactstrap'

class GoalList extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      showAddControls: false,
      title: '',
    }
    this.toggleAdd = this.toggleAdd.bind(this)
    this.handleTitleChange = this.handleTitleChange.bind(this)
  }
  toggleAdd() {
    this.setState(state => ({ showAddControls: !state.showAddControls }))
  }
  handleTitleChange(e) {
    this.setState({ title: e.target.value })
  }
  render () {
    let { goals, loading, error, startTimer, timerRefetch, setGoalDone, addWeeklyGoal, addMonthlyGoal } = this.props
    return (
      <React.Fragment>
        <div className="listTitle" onClick={this.toggleAdd}>Goals</div>
        <Container>
          {loading && <Spinner color="primary" />}
          {!loading && console.log('got data', goals)}
          {error && <div>Goals: {error.message}</div>}
          <Collapse isOpen={this.state.showAddControls}>
            <Form>
              <Input type="text"
                value={this.state.title}
                size={'50'}
                id="newGoalTitle"
                placeholder="goal description"
                onChange={this.handleTitleChange}/>
              <Button key={'add'} onClick={() => {
                addMonthlyGoal.mutation({
                  variables: {
                    title: this.state.title,
                  }
                })
                  .then(() => {
                    this.setState({ showAddControls: false })
                  })
              }} size="sm" color="primary">
                Add monthly goal</Button>
            </Form>
          </Collapse>
          {!loading && !error && 
            goals
              .map((g) => <Goal key={g.title} startTimer={startTimer} setGoalDone={setGoalDone} 
                timerRefetch={timerRefetch} addWeeklyGoal={addWeeklyGoal}
                goal={g}
                />
              )
          }
        </Container>
      </React.Fragment>
    )
  }
}

export default GoalList
