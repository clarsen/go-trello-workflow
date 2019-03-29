import auth0 from 'auth0-js'

const AUTH0_CLIENT_ID = 'PML4dRWOkZuPvU1THvcEN56a9k8JZbXh'
const AUTH0_DOMAIN = 'clarsen.auth0.com'

class Auth {
  constructor() {
    console.log('creating Auth()')
    if (!process.browser) {
      return
    }

    let redirectUri = `${window.location.protocol}//${window.location.host}/callback`
    this.auth0 = new auth0.WebAuth({
      domain: AUTH0_DOMAIN,
      clientID: AUTH0_CLIENT_ID,
      redirectUri,
      responseType: 'token id_token',
      scope: 'openid profile email'
    })

    this.state = {
      expiresAt: 0,
    }
    this.authFlag = 'workflow.isLoggedIn'
    this.authResult = 'workflow.authResult'

    this.login = this.login.bind(this)
    this.logout = this.logout.bind(this)
    this.handleAuthentication = this.handleAuthentication.bind(this)
    this.setSession = this.setSession.bind(this)
    this.silentAuth = this.silentAuth.bind(this)
    this.isAuthenticated = this.isAuthenticated.bind(this)
    this.rehydrate()
  }
  rehydrate() {
    let res = localStorage.getItem(this.authResult)
    console.log('rehydrate authResult', res)
    if (res) {
      let ar = JSON.parse(res)
      this.idToken = ar.idToken
      this.idTokenPayload = ar.idTokenPayload
    }
  }
  login() {
    this.auth0.authorize()
  }

  getIdToken() {
    return this.idToken
  }
  getName() {
    return this.idTokenPayload.name
  }
  getEmail() {
    return this.idTokenPayload.email
  }
  getPicture() {
    return this.idTokenPayload.picture
  }
  handleAuthentication() {
    return new Promise((resolve, reject) => {
      this.auth0.parseHash((err, authResult) => {
        if (err) return reject(err)
        if (!authResult || !authResult.idToken) {
          return reject(err)
        }
        this.setSession(authResult)
        resolve()
      })
    })
  }

  setSession(authResult) {
    console.log('setSession authResult=', authResult)
    this.idToken = authResult.idToken
    this.idTokenPayload = authResult.idTokenPayload
    localStorage.setItem(this.authResult, JSON.stringify(authResult))
    localStorage.setItem(this.authFlag, JSON.stringify(true))
    console.log(this.idToken)
  }

  logout() {
    let redirectUri = `${window.location.protocol}//${window.location.host}/`
    localStorage.setItem(this.authFlag, JSON.stringify(false))
    localStorage.removeItem(this.authResult)
    this.auth0.logout({
      returnTo: redirectUri,
      clientID: AUTH0_CLIENT_ID,
    })
  }

  silentAuth() {
    console.log('silentAuth')
    if(this.isAuthenticated()) {
      console.log('  isAuthenticated')
      return new Promise((resolve, reject) => {
        console.log('  silentAuth promise calling')
        this.auth0.checkSession({}, (err, authResult) => {
          console.log('  checked session got', err, authResult)
          if (err) {
            localStorage.removeItem(this.authFlag)
            return reject(err)
          }
          this.setSession(authResult)
          resolve()
        })
      })
    }
  }

  isAuthenticated() {
    if (!process.browser) {
      return false
    }
    console.log('isAuthenticated Auth=', this)
    let res = JSON.parse(localStorage.getItem(this.authFlag))
    console.log('isAuthenticated Auth=', this, 'localstorage =', res, 'idToken =', this.getIdToken())
    return res && this.getIdToken() !== undefined // only if auth *and* idtoken available
  }
}

var _auth = null

const auth = () => {
  if (_auth === null) {
    _auth = new Auth()
  }
  return _auth
}
export default auth
