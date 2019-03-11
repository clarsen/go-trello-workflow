import withApollo from 'next-with-apollo'
import ApolloClient, { InMemoryCache }  from 'apollo-boost'
import { ENDPOINTS } from './api'
import auth from './auth0'

export default withApollo(({ ctx, headers, initialState }) => (
  new ApolloClient({
    uri: ENDPOINTS['go']['private_gql'],
    cache: new InMemoryCache().restore(initialState || {}),
    request: operation => {
      operation.setContext(context => ({
        headers: {
          ...context.headers,
          Authorization: `Bearer ${auth().getIdToken()}`,
        },
      }))
    }
  })
))
