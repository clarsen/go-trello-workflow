import React from 'react'
import {
  Collapse,
  Input,
  Spinner,
  Table,
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
        <td>{dur_for_day[day]>0 && `${numeral(dur_for_day[day]/1000.0).format('00:00:00')}`}</td>
      )
    }

    return (
      <React.Fragment>
        <tr>
          <td>{entry.detail}</td>
          <td>{
            `${numeral(tot_ms/1000.0).format('00:00:00')}`
          }</td>
        </tr>
        <Collapse isOpen={this.props.showDetails}>
          <tr>
            <td>Su</td>
            <td>M</td>
            <td>Tu</td>
            <td>W</td>
            <td>Th</td>
            <td>F</td>
            <td>Sa</td>
          </tr>
          <tr>
            {cols}
          </tr>
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
        <tr className="projectItem">
          <th scope="row" onClick={this.toggle}>{project.title}</th>
          <td>{
            `${numeral(tot_ms/1000.0).format('00:00:00')}`
          }</td>
          <td>
            <Table dark>
              {project.entries.map((e) =>
                <ProjectDetailEntry entry={e} showDetails={this.state.showDetails}/>
              )}
            </Table>
          </td>
        </tr>
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
            <Table dark striped>
              <thead>
                <th>Project</th>
                <th>Total</th>
                <th>Detail/Duration (HH:MM:SS)</th>
              </thead>
              <tbody>
                {data.projects.map((p) => <ProjectItem project={p} />)}
              </tbody>
            </Table>
            }
            <style global jsx>{`
              .weekSelect {
                display: inline;
                width: 4em;
              }
            `}</style>
          </React.Fragment>
        }
      </QueryContainer>
    )
  }
}

export default ProjectReport
