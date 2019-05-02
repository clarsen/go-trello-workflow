/**
 * Layout component that queries for data
 * with Gatsby's StaticQuery component
 *
 * See: https://www.gatsbyjs.org/docs/static-query/
 */

import React from "react"
import PropTypes from "prop-types"
import { Link, StaticQuery, graphql } from "gatsby"
// import { css } from "@emotion/core"

import { ApolloProvider } from 'react-apollo'
import ApolloClient, { InMemoryCache }  from 'apollo-boost'
import TopWrap from '../TopWrap'
import { Provider as AlertProvider } from 'react-alert'
import AlertTemplate from 'react-alert-template-basic'
import { ENDPOINTS } from '../../lib/api'
import auth from '../../lib/auth0'
import 'cross-fetch/polyfill'

// main site style
import 'bootstrap/dist/css/bootstrap.min.css'

// optional cofiguration
const alertOptions = {
  timeout: 5000,
}

const apollo = new ApolloClient({
  uri: ENDPOINTS['go']['private_gql'],
  fetch,
  // cache: new InMemoryCache().restore(initialState || {}),
  request: operation => {
    operation.setContext(context => ({
      headers: {
        ...context.headers,
        Authorization: `Bearer ${auth.instance().getIdToken()}`,
      },
    }))
  }
})

const Layout = ({ children }) => (
  <StaticQuery
    query={graphql`
      query SiteTitleQuery {
        site {
          siteMetadata {
            title
          }
        }
      }
    `}
    render={data => (
      <ApolloProvider client={apollo}>
        <AlertProvider template={AlertTemplate} {...alertOptions}>
          <TopWrap>
            {children}
          </TopWrap>
        </AlertProvider>
      </ApolloProvider>
    )}
  />
)

Layout.propTypes = {
  children: PropTypes.node.isRequired,
}

export default Layout
