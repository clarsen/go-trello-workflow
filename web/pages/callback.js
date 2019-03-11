import React from 'react'
import Router from 'next/router'
import auth from '../lib/auth0'

class CallbackPage extends React.Component {
  async componentDidMount() {
    await auth().handleAuthentication()
    Router.push('/')
  }
  render () {
    return (
      <div>This is authentication page</div>
    )
  }
}

export default CallbackPage
