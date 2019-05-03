import { navigate } from 'gatsby'
import React from 'react'
import auth from '../lib/auth0'

interface ICustomInputProps {
  path: string,
}

class LoginCallback extends React.Component<ICustomInputProps> {
  public async componentDidMount() {
    console.log('LoginCallback componentDidMount')
    try {
      await auth.instance().handleAuthentication()
      navigate('/')
    } catch (err) {
      console.log(err)
    }
  }
  public render() {
    const authFlag = localStorage.getItem(auth.instance().authFlag)
    return (
      <div>This is authentication page.
        auth = {`${auth.instance()}`},
        isAuthenticated = {`${auth.instance().isAuthenticated()}`}.
        getIdToken = {`${auth.instance().getIdToken()}`},
        authFlag = {`${auth.instance().authFlag}`},
        authFlagValue in storage = {`${authFlag}`}
        </div>
    )
  }
}

export default LoginCallback
