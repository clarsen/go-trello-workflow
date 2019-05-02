import { navigate } from 'gatsby'
import React from 'react'
import auth from '../lib/auth0'

interface ICustomInputProps {
  path: string,
}

class Login extends React.Component<ICustomInputProps> {
  public async componentDidMount() {
    if (auth.instance().isAuthenticated()) {
      navigate('/')
      return null
    }
    await auth.instance().login()
  }
  public render() {
    return (
      <div>Logging in...</div>
    )
  }
}

export default Login
