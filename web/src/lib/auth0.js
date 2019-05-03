import auth0 from 'auth0-js'

const AUTH0_CLIENT_ID = 'PML4dRWOkZuPvU1THvcEN56a9k8JZbXh'
const AUTH0_DOMAIN = 'clarsen.auth0.com'
import moment from 'moment'

class Auth {
  constructor() {
    console.log('creating Auth()')
    if (!process.browser) {
      return
    }
    // think the Safari browser will redirect to 'callback/' if we don't have the trailing slash in the URL
    let redirectUri = `${window.location.protocol}//${window.location.host}/callback/`
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
    this.authExpiresAfter = 'workflow.authExpiresAfter'

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
    console.log('auth().login - calling auth0 authorize')
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
      console.log('calling auth0.parseHash, location=', window.location, ' hash=', window.location.hash)
      if (window.location.hash === "") {
        return reject(new Error("No location hash"))
      }
      this.auth0.parseHash((err, authResult) => {
        if (err) {
          console.log(err)
          return reject(err)
        }
        if (!authResult || !authResult.idToken) {
          console.log('authResult =', authResult)
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
    let expAfter = authResult.expiresIn + moment().unix()
    console.log(`expires after ${expAfter} (${authResult.expiresIn})`)
    localStorage.setItem(this.authExpiresAfter, JSON.stringify(expAfter))
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
      let expTime = JSON.parse(localStorage.getItem(this.authExpiresAfter))
      if (expTime > moment().unix()) {
        console.log(` still valid until ${expTime}`)
        return
      }
      return new Promise((resolve, reject) => {
        console.log('  silentAuth promise calling')
        this.auth0.checkSession({}, (err, authResult) => {
          console.log('  checked session got', err, authResult)
          if (err) {
            localStorage.removeItem(this.authExpiresAfter)
            localStorage.removeItem(this.authFlag)
            localStorage.removeItem(this.authResult)
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
    let authFlag = localStorage.getItem(this.authFlag)
    console.log(`isAuthenticated ${this.authFlag} authFlag=`, authFlag)
    let res = JSON.parse(authFlag)
    console.log('isAuthenticated Auth=', this, 'authFlag parsed =', res, 'idToken =', this.getIdToken())
    return res && this.getIdToken() !== undefined // only if auth *and* idtoken available
  }
}


const authSingleton = (() => {
  var _auth = null
  return {
    instance: () => {
      if (_auth === null) {
        _auth = new Auth()
      }
      return _auth
    }
  }
})()

export default authSingleton
