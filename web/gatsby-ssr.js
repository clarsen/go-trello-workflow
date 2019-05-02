/**
 * Implement Gatsby's SSR (Server Side Rendering) APIs in this file.
 *
 * See: https://www.gatsbyjs.org/docs/ssr-apis/
 */

import { ApolloProvider } from 'react-apollo'
import ApolloClient from 'apollo-boost'
import { ENDPOINTS } from './src/lib/api'
import 'cross-fetch/polyfill'

const client = new ApolloClient({
  uri: 'http://localhost:8080/api/gql',
  fetch,
})

export const wrapRootElement = ({ element }) => (
  <ApolloProvider client={client}>{element}</ApolloProvider>
)
