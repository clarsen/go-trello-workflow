import React from 'react'
import {
  Container,
  Row,
  Col,
  Spinner,
  Input,
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
    let tot_ms = project.entries.reduce((acc, e) =>
        acc + e.entries.reduce((acc, de) =>
          acc + de.duration_ms, 0), 0)

    return (
      <React.Fragment>
        <Row>
          <Col lg={2}>{project.title}</Col>
          <Col lg={1}>{
            `${numeral(tot_ms/1000.0).format('00:00:00')}`
          }</Col>
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
    let nowGrace = moment().subtract(3,'days')
    let nowGraceMonth = moment().subtract(5,'days')
    let monthNext = nowGraceMonth.month()+1
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
