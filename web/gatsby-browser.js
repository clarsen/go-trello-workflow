/**
 * Implement Gatsby's Browser APIs in this file.
 *
 * See: https://www.gatsbyjs.org/docs/browser-apis/
 */

// You can delete this file if you're not using it

import ApolloClient from 'apollo-boost'
import { ApolloProvider } from 'react-apollo'
import React from 'react'
import { ENDPOINTS } from './src/lib/api'
import 'cross-fetch/polyfill'

const client = new ApolloClient({
  uri: ENDPOINTS['go']['private_gql'],
  fetch
  // cache: new InMemoryCache().restore(initialState || {})
})

export const wrapRootElement = ({ element }) => (
  <ApolloProvider client={client}>{element}</ApolloProvider>
)
