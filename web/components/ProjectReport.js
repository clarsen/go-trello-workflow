import React from 'react'
import {
  Container,
  Row,
  Col,
  Spinner,
} from 'reactstrap'
import numeral from 'numeral'

class ProjectDetailEntry extends React.Component {
  constructor (props) {
    super(props)
  }
  render() {
    let { entry } = this.props
    let tot_ms = entry.entries.reduce((acc, entr) => acc + entr.duration_ms, 0)

    return (
      <React.Fragment>
        <Row>
          <Col>{entry.detail}</Col>
          <Col>{
            `${numeral(tot_ms/1000.0).format('00:00:00')}`
          }</Col>
        </Row>
      </React.Fragment>
    )
  }
}

class ProjectItem extends React.Component {
  constructor (props) {
    super(props)
  }
  render() {
    let { project } = this.props
    return (
      <React.Fragment>
        <Row>
          <Col lg={2}>{project.title}</Col>
          <Col>
            {project.entries.map((e) => <ProjectDetailEntry entry={e} />)}
          </Col>
        </Row>
      </React.Fragment>
    )
  }
}

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
          <Row className="rowHeader"><Col>Project</Col><Col>Detail</Col><Col>Duration (HH:MM:SS)</Col></Row>
          {data.projects.map((p) => <ProjectItem project={p} />)}
        </Container>
        }
        <style global jsx>{`
          .rowHeader {
            text-decoration: underline
          }
        `}</style>
      </React.Fragment>
    )
  }
}

export default ProjectReport
