import App, { Container } from 'next/app'
import { ApolloProvider } from 'react-apollo'
import withApollo from '../lib/withApollo'
import TopWrap from '../components/TopWrap'
import 'bootstrap/dist/css/bootstrap.min.css'
import { Provider as AlertProvider } from 'react-alert'
import AlertTemplate from 'react-alert-template-basic'

// optional cofiguration
const alertOptions = {
  timeout: 5000,
}

class MyApp extends App {
  render() {
    const { Component, pageProps, apollo } = this.props

    return (
      <Container>
        <ApolloProvider client={apollo}>
          <AlertProvider template={AlertTemplate} {...alertOptions}>
            <TopWrap>
              <Component {...pageProps} />
            </TopWrap>
          </AlertProvider>
        </ApolloProvider>
      </Container>
    )
  }
}

export default withApollo(MyApp)
