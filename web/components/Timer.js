import React from 'react'
import {
  Button,
  Row,
  Col,
} from 'reactstrap'

class Timer extends React.Component {
  constructor (props) {
    super(props)
  }
  render () {
    let { activeTimer, stopTimer, timerRefetch } = this.props
    return (
      <React.Fragment>
        {!activeTimer && <div>No timer active</div>}
        {activeTimer &&
          <Row key={'row0'+activeTimer.title}>
            <Col xs={12} lg={12} key={'1'+activeTimer.title}>
              Currently active: {activeTimer.title}
              <Button outline color='primary' size='sm' onClick={()=>{
                stopTimer.mutation({
                  variables: {
                    timerID: activeTimer.id,
                  }
                })
                  .then(() => timerRefetch())
              }}>Stop</Button>
            </Col>
          </Row>
        }
      </React.Fragment>

    )
  }
}

export default Timer
