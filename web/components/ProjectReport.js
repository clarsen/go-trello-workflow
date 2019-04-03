import React from 'react'
import {
  Col,
  Collapse,
  Container,
  Input,
  Row,
  Spinner,
} from 'reactstrap'
import { FaSync } from 'react-icons/fa'

import numeral from 'numeral'
import moment from 'moment'
import { Query } from 'react-apollo'
import { adopt } from 'react-adopt'
import ApolloClient, { InMemoryCache }  from 'apollo-boost'

import {
  ProjectReportQuery,
} from '../lib/timereport_graphql'
import { ENDPOINTS } from '../lib/api'
import auth from '../lib/auth0'


const pythonGraphqlClient = new ApolloClient({
  uri: ENDPOINTS['python']['private_gql'],
  cache: new InMemoryCache().restore({}),
  request: operation => {
    operation.setContext(context => ({
      headers: {
        ...context.headers,
        Authorization: `Bearer ${auth().getIdToken()}`,
      },
    }))
  },
})

const QueryContainer = adopt({
  projectReportQuery: ({ render, year, week }) => (
    <Query client={pythonGraphqlClient} query={ProjectReportQuery} ssr={false} variables={{ year, week }}>
      {render}
    </Query>
  ),
})

class ProjectDetailEntry extends React.Component {
  constructor (props) {
    super(props)
  }
  render() {
    let { entry, showDetails } = this.props
    let tot_ms = entry.entries
      .reduce((acc, entr) => acc + entr.duration_ms, 0)
    let dur_for_day = new Map()
    let cols = []
    for (let day = 0; day < 7; day++) {
      dur_for_day[day] = entry.entries
        .filter(e => moment.unix(e.start).day() === day )
        .reduce((acc, e) => acc + e.duration_ms, 0)
      cols.push(
        <Col lg={1}>{dur_for_day[day]>0 && `${numeral(dur_for_day[day]/1000.0).format('00:00:00')}`}</Col>
      )
    }

    return (
      <React.Fragment>
        <Row className="projectDetailEntry">
          <Col>{entry.detail}</Col>
          <Col>{
            `${numeral(tot_ms/1000.0).format('00:00:00')}`
          }</Col>
        </Row>
        <Collapse isOpen={this.props.showDetails}>
          <Row className="dayHoursDetail">
            <Col lg={1}>Su</Col>
            <Col lg={1}>M</Col>
            <Col lg={1}>Tu</Col>
            <Col lg={1}>W</Col>
            <Col lg={1}>Th</Col>
            <Col lg={1}>F</Col>
            <Col lg={1}>Sa</Col>
          </Row>
          <Row className="dayHoursDetail">
            {cols}
          </Row>
        </Collapse>
      </React.Fragment>
    )
  }
}

class ProjectItem extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      showDetails: false
    }
    this.toggle = this.toggle.bind(this)
  }
  toggle() {
    this.setState(state => ({ showDetails: !state.showDetails }))
  }
  render() {
    let { project } = this.props
    let tot_ms = project.entries.reduce((acc, e) =>
      acc + e.entries.reduce((acc, de) =>
      acc + de.duration_ms, 0), 0)

    return (
      <React.Fragment>
        <Row className="projectItem">
          <Col lg={2} onClick={this.toggle}>{project.title}</Col>
          <Col lg={1}>{
            `${numeral(tot_ms/1000.0).format('00:00:00')}`
          }</Col>
          <Col>
            {project.entries.map((e) =>
              <ProjectDetailEntry entry={e} showDetails={this.state.showDetails}/>
            )}
          </Col>
        </Row>
      </React.Fragment>
    )
  }
}

class ProjectReport extends React.Component {
  constructor (props) {
    super(props)
    let now = moment()
    this.state = {
      week: now.isoWeek()
    }
    this.changeWeek = this.changeWeek.bind(this)
  }
  changeWeek(e) {
    this.setState({ week: parseInt(e.target.value) })
  }
  render() {
    let now = moment()
    let nowGraceMonth = moment().subtract(5,'days')
    return (
      <QueryContainer year={now.year()} week={this.state.week} month={3}>
        {({
          projectReportQuery: { loading, data, error, refetch: projectReportRefetch },
        }) =>
          <React.Fragment>
            {'Week '}<Input className="weekSelect" type="select" id="week" value={this.state.week} onChange={this.changeWeek}>
              <option>10</option>
              <option>11</option>
              <option>12</option>
              <option>13</option>
              <option>14</option>
              <option>15</option>
              <option>16</option>
              <option>17</option>
              <option>18</option>
            </Input>{' '}
            <FaSync size={25} onClick={() => {
              projectReportRefetch()
            }} />
            {loading && <Spinner color="primary" />}
            {!loading && console.log('got data', data)}
            {(!loading && !error) &&
            <Container>
              <Row className="rowHeader">
                <Col lg={2}>Project</Col>
                <Col lg={1}>Total</Col>
                <Col>Detail</Col>
                <Col>Duration (HH:MM:SS)</Col>
              </Row>
              {data.projects.map((p) => <ProjectItem project={p} />)}
            </Container>
            }
            <style global jsx>{`
              .weekSelect {
                display: inline;
                width: 3em;
              }
              .projectDetailEntry {
                background-color: #222;
              }
              .projectItem {
                background-color: #222;
              }
              .dayHoursDetail {
                background-color: #000;
              }
              .rowHeader {
                text-decoration: underline
              }
            `}</style>
          </React.Fragment>
        }
      </QueryContainer>
    )
  }
}

export default ProjectReport
