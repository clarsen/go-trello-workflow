import React from 'react'
import {
  Collapse,
  Input,
  Navbar,
  Nav,
  NavItem,
  NavLink,
  Spinner,
  Table,
  TabContent,
  TabPane
} from 'reactstrap'
import { FaSync } from 'react-icons/fa'
import classnames from 'classnames'
import 'cross-fetch/polyfill'

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
  fetch,
  cache: new InMemoryCache().restore({}),
  request: operation => {
    operation.setContext(context => ({
      headers: {
        ...context.headers,
        Authorization: `Bearer ${auth.instance().getIdToken()}`,
      },
    }))
  },
})

const QueryContainer = adopt({
  projectWeeklyReportQuery: ({ render, year, week }) => (
    <Query client={pythonGraphqlClient} query={ProjectReportQuery} ssr={false} variables={{ year, week }}>
      {render}
    </Query>
  ),
  projectMonthlyReportQuery: ({ render, year, month }) => (
    <Query client={pythonGraphqlClient} query={ProjectReportQuery} ssr={false} variables={{ year, month }}>
      {render}
    </Query>
  ),
  projectYearlyReportQuery: ({ render, year }) => (
    <Query client={pythonGraphqlClient} query={ProjectReportQuery} ssr={false} variables={{ year }}>
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
    let rollup = project.entries.reduce((dur_for_day, entry) => {
        for (let day = 0; day < 7; day++) {
          dur_for_day[day] += entry.entries
            .filter(e => moment.unix(e.start).day() === day )
            .reduce((acc, e) => acc + e.duration_ms, 0)
        }
        return dur_for_day
      }, 
      {0: 0, 1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0}
    )

    let rollupCols =[]
    for (let day = 0; day < 7; day++) {
      rollupCols.push(
        <td>{rollup[day]>0 && `${numeral(rollup[day]/1000.0).format('00:00:00')}`}</td>
      )
    }

    return (
      <React.Fragment>
        <tr className="projectItem">
          <th scope="row" onClick={this.toggle}>{project.title}</th>
          <td>{
            `${numeral(tot_ms/1000.0).format('00:00:00')}`
          }</td>
          <td>
            <Collapse isOpen={this.state.showDetails}>
              <Table dark>
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
                  {rollupCols}
                </tr>
                </Table>
            </Collapse>
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

function jsonCopy(src) {
  return JSON.parse(JSON.stringify(src));
}

const rebucket = (projects, levels) => {
  let rebucketed = projects.reduce((newProjects, p) => {
    let comps = p.title.split(/\s*-\s*/)
    let newTitle = comps.slice(0,levels).join(' - ')
    if (newTitle in newProjects)   {
      newProjects[newTitle].entries = newProjects[newTitle].entries.concat(jsonCopy(p.entries))
    } else {
      newProjects[newTitle] = jsonCopy(p)
      newProjects[newTitle].title = newTitle
    }
    return newProjects
  }, {})
  let newProjects = []
  for (var key in rebucketed) {
    newProjects.push(rebucketed[key])
  }
  return newProjects
}

class ProjectReport extends React.Component {
  constructor (props) {
    super(props)
    let now = moment()
    this.state = {
      week: now.isoWeek(),
      month: now.month() + 1,
      year: now.year(),
      activeTab: 'weekReport',
    }
    this.changeWeek = this.changeWeek.bind(this)
    this.changeMonth = this.changeMonth.bind(this)
    this.changeYear = this.changeYear.bind(this)
    this.switchTab = this.switchTab.bind(this)
  }

  switchTab (tab) {
    if (this.state.activeTab !== tab) {
      this.setState({
        activeTab: tab
      })
    }
  }

  changeWeek(e) {
    this.setState({ week: parseInt(e.target.value) })
  }

  changeMonth(e) {
    this.setState({ month: parseInt(e.target.value) })
  }

  changeYear(e) {
    this.setState({ year: parseInt(e.target.value) })
  }

  render() {
    let now = moment()
    let nowGraceMonth = moment().subtract(5,'days')
    return (
      <QueryContainer year={now.year()} week={this.state.week} month={this.state.month}>
        {({
          projectWeeklyReportQuery: { loading, data, error, refetch: projectWeeklyReportRefetch },
          projectMonthlyReportQuery: { loading: loadingMonthly, data: dataMonthly, error: errorMonthly, refetch: projectMonthlyReportRefetch },
          projectYearlyReportQuery: { loading: loadingYearly, data: dataYearly, error: errorYearly, refetch: projectYearlyReportRefetch },
        }) =>
          <React.Fragment>
            <Navbar expand="lg">
              <Nav className="mr-auto" tabs>
                <NavItem>
                  <NavLink className={classnames({ active: this.state.activeTab === 'weekReport'})} onClick={()=> this.switchTab('weekReport')}>
                    Weekly
                  </NavLink>
                </NavItem>
                <NavItem>
                  <NavLink className={classnames({ active: this.state.activeTab === 'monthReport'})} onClick={()=> this.switchTab('monthReport')}>
                    Monthly
                  </NavLink>
                </NavItem>
                <NavItem>
                  <NavLink className={classnames({ active: this.state.activeTab === 'yearReport'})} onClick={()=> this.switchTab('yearReport')}>
                    Yearly
                  </NavLink>
                </NavItem>
              </Nav>
            </Navbar>
            <TabContent activeTab={this.state.activeTab}>
              <TabPane tabId="weekReport">
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
                  projectWeeklyReportRefetch()
                }} />
                <div className="sectionName">Project summary:</div>
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
                <div className="sectionName">Project group summary:</div>
                {(!loading && !error) &&
                <Table dark striped>
                  <thead>
                    <th>Project</th>
                    <th>Total</th>
                    <th>Detail/Duration (HH:MM:SS)</th>
                  </thead>
                  <tbody>
                    {rebucket(data.projects, 2).map((p) => <ProjectItem project={p} />)}
                  </tbody>
                </Table>
                }
                <style global jsx>{`
                  .weekSelect {
                    display: inline;
                    width: 4em;
                  }
                  .sectionName {
                    font-size: 1.3em;
                    font-weight: bold;
                    text-decoration: underline;
                  }
                `}</style>
              </TabPane>
              <TabPane tabId="monthReport">
                {'Month '}<Input className="monthSelect" type="select" id="month" value={this.state.month} onChange={this.changeMonth}>
                    <option>1</option>
                    <option>2</option>
                    <option>3</option>
                    <option>4</option>
                    <option>5</option>
                </Input>{' '}
                <FaSync size={25} onClick={() => {
                  projectMonthlyReportRefetch()
                }} />
                <div className="sectionName">Project summary:</div>
                {loadingMonthly && <Spinner color="primary" />}
                {!loadingMonthly && console.log('got data', dataMonthly)}
                {(!loadingMonthly && !errorMonthly) &&
                <Table dark striped>
                  <thead>
                    <tr>
                      <th>Project</th>
                      <th>Total</th>
                      <th>Detail/Duration (HH:MM:SS)</th>
                    </tr>
                  </thead>
                  <tbody>
                    {dataMonthly.projects.map((p) => <ProjectItem project={p} />)}
                  </tbody>
                </Table>
                }
                <div className="sectionName">Project group summary:</div>
                {(!loadingMonthly && !errorMonthly) &&
                <Table dark striped>
                  <thead>
                    <tr>
                      <th>Project</th>
                      <th>Total</th>
                      <th>Detail/Duration (HH:MM:SS)</th>
                    </tr>
                  </thead>
                  <tbody>
                    {rebucket(dataMonthly.projects, 2).map((p) => <ProjectItem project={p} />)}
                  </tbody>
                </Table>
                }
                <style global jsx>{`
                  .monthSelect {
                    display: inline;
                    width: 4em;
                  }
                  .sectionName {
                    font-size: 1.3em;
                    font-weight: bold;
                    text-decoration: underline;
                  }
                `}</style>
              </TabPane>
              <TabPane tabId="yearReport">
                {'Year '}<Input className="yearSelect" type="select" id="year" value={this.state.year} onChange={this.changeYear}>
                    <option>2018</option>
                    <option>2019</option>
                </Input>{' '}
                <FaSync size={25} onClick={() => {
                  projectYearlyReportRefetch()
                }} />
                <div className="sectionName">Project summary:</div>
                {loadingYearly && <Spinner color="primary" />}
                {!loadingYearly && console.log('got data', dataYearly)}
                {(!loadingYearly && !errorYearly) &&
                <Table dark striped>
                  <thead>
                    <tr>
                      <th>Project</th>
                      <th>Total</th>
                      <th>Detail/Duration (HH:MM:SS)</th>
                    </tr>
                  </thead>
                  <tbody>
                    {dataYearly.projects.map((p) => <ProjectItem project={p} />)}
                  </tbody>
                </Table>
                }
                <div className="sectionName">Project group summary:</div>
                {(!loadingYearly && !errorYearly) &&
                <Table dark striped>
                  <thead>
                    <tr>
                      <th>Project</th>
                      <th>Total</th>
                      <th>Detail/Duration (HH:MM:SS)</th>
                    </tr>
                  </thead>
                  <tbody>
                    {rebucket(dataYearly.projects, 2).map((p) => <ProjectItem project={p} />)}
                  </tbody>
                </Table>
                }
                <style global jsx>{`
                  .yearSelect {
                    display: inline;
                    width: 5em;
                  }
                  .sectionName {
                    font-size: 1.3em;
                    font-weight: bold;
                    text-decoration: underline;
                  }
                `}</style>
              </TabPane>
            </TabContent>
          </React.Fragment>
        }
      </QueryContainer>
    )
  }
}

export default ProjectReport
