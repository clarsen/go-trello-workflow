import { navigate } from 'gatsby'
import React from 'react'
import auth from '../lib/auth0'

interface ICustomInputProps {
  path: string,
}

class LoginCallback extends React.Component<ICustomInputProps> {
  public async componentDidMount() {
    await auth().handleAuthentication()
    navigate('/')
  }
  public render() {
    return (
      <div>This is authentication page</div>
    )
  }
}

export default LoginCallback
