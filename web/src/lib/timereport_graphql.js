import gql from 'graphql-tag'


export const TimeReportQuery = gql`
  query report {
    default_report
  }
`

export const ProjectReportQuery = gql`
  query projects($year: Int, $month: Int, $week: Int) {
    projects(year: $year, month: $month, week: $week) {
      title
      pid
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
