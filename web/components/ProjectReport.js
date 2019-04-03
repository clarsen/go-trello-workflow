import React from 'react'
import {
  Container,
  Row,
  Col,
  Spinner,
} from 'reactstrap'

class ProjectReport extends React.Component {
  constructor (props) {
    super(props)
  }
  render() {
    let { loading, error, data } = this.props
    return (
      <React.Fragment>
      {loading && <Spinner color="primary" />}
      {!loading && console.log('got data', data)}
      {(!loading && !error) &&
        <Container>
          {data.projects.map((p) => <Row>
            <Col lg={2}>{p.title}</Col>
            <Col>
            {p.entries.map((e) => <Row>
              <Col>{e.detail}</Col>
              <Col>{
                `${e.entries.reduce((acc, entr) => acc + entr.duration_ms, 0)/60.0/1000.0}`
              }</Col>
              </Row>
            )}
            </Col>
            </Row>)
          }
        </Container>
      }
      </React.Fragment>
    )
  }
}

export default ProjectReport
