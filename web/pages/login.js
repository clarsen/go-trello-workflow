import React from 'react'
import {
  Button,
} from 'reactstrap'
import auth from '../lib/auth0'

class LoginPage extends React.Component {
  render () {
    return (
      <Button onClick={() => { console.log('Login'); auth().login() }}>You must log in</Button>
    )
  }
}

export default LoginPage
