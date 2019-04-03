import gql from 'graphql-tag'


export const TimeReportQuery = gql`
  query report {
    default_report
  }
`

export const ProjectReportQuery = gql`
  query projects($year: Int, $week: Int) {
    projects(year: $year, week: $week) {
      title
      entries {
        detail
        entries {
          id
          start
          duration_ms
        }
      }
    }
  }
`
