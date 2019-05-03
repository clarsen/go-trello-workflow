import { navigate } from 'gatsby'
import React from 'react'
import auth from '../lib/auth0'

interface ICustomInputProps {
  path: string,
}

class Login extends React.Component<ICustomInputProps> {
  public render() {
    if (auth.instance().isAuthenticated()) {
      navigate('/app')
      return null
    }
    auth.instance().login()
    return (
      <div>Logging in...</div>
    )
  }
}

export default Login
