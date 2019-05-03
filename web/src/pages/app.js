import { Router } from '@reach/router'
import { graphql } from 'gatsby'
import React from 'react'

import Layout from '../components/layout'
import LoginCallback from '../components/LoginCallback'
import IndexPage from '../components/IndexPage'
import Login from '../components/Login'

export default () => {
  return (
    <Layout>
      <Router>
        <IndexPage path='/app' />
        <Login path='/app/login' />
        <LoginCallback path='/app/callback/' />
      </Router>
    </Layout>
  )
}
