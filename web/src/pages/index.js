import { Router } from '@reach/router'
import { graphql } from 'gatsby'
import React from 'react'
import { Link } from 'gatsby'

import Layout from '../components/layout'
import LoginCallback from '../components/LoginCallback'
import IndexPage from '../components/IndexPage'
import Login from '../components/Login'

export default () => {
  return (
    <Layout>
      <Link to='/app'>
        <b> Go to app</b>
      </Link>
    </Layout>
  )
}
