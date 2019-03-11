import React from 'react'
import { Alert, Button } from 'reactstrap'

import { Query, Mutation } from 'react-apollo'
import { adopt } from 'react-adopt'
import gql from 'graphql-tag'
import { withAlert } from 'react-alert'


import NavHeader from '../components/NavHeader'
import TaskList from '../components/TaskList'
import MarkdownRenderer from 'react-markdown-renderer'


export const taskQuery = gql`
  query tasks($inBoardList: BoardList) {
    tasks(inBoardList: $inBoardList) {
      id
      title
      createdDate
      url
    }
  }
`

const generateWeeklySummaryQuery = gql`
  mutation generateWeeklySummary {
    generateWeeklySummary
  }
`

const generateWeeklySummary = ({ render }) => (
  <Mutation
    mutation={generateWeeklySummaryQuery}
  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)

const generateWeeklyReviewTemplateQuery = gql`
  mutation generateWeeklyTemplate {
    generateWeeklyReviewTemplate
  }
`
const generateWeeklyReviewTemplate = ({ render }) => (
  <Mutation
    mutation={generateWeeklyReviewTemplateQuery}
  >
    {(mutation, result) => render({ mutation, result })}
  </Mutation>
)


const weeklyVisualizationQuery = gql`
  query weeklyVisualization {
    weeklyVisualization
  }
`

/* eslint-disable react/display-name */
/* eslint-disable react/prop-types */

const QueryContainer = adopt({
  queryAll: ({ render }) => (
    <Query query={taskQuery} ssr={false} variables={{ }}>
      {render}
    </Query>
  ),
  weeklyVisualizationQuery: ({ render }) => (
    <Query query={weeklyVisualizationQuery} ssr={false} variables={{ }}>
      {render}
    </Query>
  ),
  generateWeeklySummary,
  generateWeeklyReviewTemplate,
})

class IndexPage extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      data: null
    }
  }
  render () {
    let { alert } = this.props
    return (
      <QueryContainer>
        {({
          queryAll: { loading: loadingAll, data: allTasks },
          weeklyVisualizationQuery: { loading : weeklyLoading, data: weeklyVisualizationData },
          generateWeeklySummary,
          generateWeeklyReviewTemplate,
        }) =>
          <React.Fragment>
            <NavHeader/>
            <Button color='primary' onClick={() => {
              generateWeeklySummary.mutation()
            }}>
                Generate Weekly Summary
            </Button>
            <Button color='primary' onClick={() => {
              generateWeeklyReviewTemplate.mutation()
                .catch((e) => {
                  console.log('error', e)
                  alert.show(e.message)
                })
            }}>
                Generate Weekly Review Template
            </Button>
            {generateWeeklyReviewTemplate.errors &&
              <Alert color="warning">
              {generateWeeklyReviewTemplate.errors.map((e) => e.message)}
              </Alert>
            }
            {loadingAll && <div>loading...</div>}
            {!loadingAll && console.log('got data', allTasks)}
            {!loadingAll && <TaskList tasks={allTasks.tasks}/>}

            {weeklyLoading && <div>loading...</div>}
            {!weeklyLoading &&
              <MarkdownRenderer markdown={weeklyVisualizationData.weeklyVisualization} />}
          </React.Fragment>
        }
      </QueryContainer>
    )
  }
}

export default withAlert()(IndexPage)
