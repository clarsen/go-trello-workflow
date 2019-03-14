import React from 'react'
import Head from 'next/head'
import {
  Navbar,
  Nav,
  NavItem,
  NavLink,
} from 'reactstrap'
import { withRouter } from 'next/router'
import auth from '../lib/auth0'

class Header extends React.Component {
  constructor(props) {
    super(props)
  }
  render () {
    // let { router : { pathname }, queue } = this.props
    return (
      <React.Fragment>
        <Head>
          <meta charSet="UTF-8" />
          <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no" />
          <style jsx>{`
            header {
              margin-bottom: 25px;
            }
            a {
              font-size: 14px;
              margin-right: 15px;
              text-decoration: none;
            }
            .is-active {
              text-decoration: underline;
            }
          `}</style>
        </Head>
        <Navbar className="sticky-top" expand="lg">
          <Nav className="mr-auto" navbar>
            {auth().isAuthenticated()
              ? <NavItem onClick={() => { console.log('Logout'); auth().logout() }}>
                <NavLink to="#" className="nav-link">
                  Logout</NavLink>
              </NavItem>
              : <NavItem onClick={() => { console.log('Login'); auth().login() }}>
                <NavLink to="#" className="nav-link">
                  Login</NavLink>
              </NavItem>
            }
          </Nav>
        </Navbar>
      </React.Fragment>

    )
  }
}

export default withRouter(Header)
