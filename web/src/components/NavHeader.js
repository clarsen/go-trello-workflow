import React from 'react'
import Head from 'next/head'
import {
  Navbar,
  Nav,
  NavItem,
  NavLink,
} from 'reactstrap'
import classnames from 'classnames'

import { withRouter } from 'next/router'
import auth from '../lib/auth0'

class Header extends React.Component {
  constructor(props) {
    super(props)
  }
  render () {
    // let { router : { pathname }, queue } = this.props
    let { switchTab, activeTab } = this.props
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
        <Navbar expand="lg">
          <Nav className="mr-auto" tabs>
            <NavItem>
              <NavLink className={classnames({ active: activeTab === 'board'})} onClick={()=> switchTab('board')}>
                Workboard
              </NavLink>
            </NavItem>
            <NavItem>
              <NavLink className={classnames({ active: activeTab === 'periodicBoard'})} onClick={()=> switchTab('periodicBoard')}>
                Periodic
              </NavLink>
            </NavItem>
            <NavItem>
              <NavLink className={classnames({ active: activeTab === 'oldTasksReview'})} onClick={()=> switchTab('oldTasksReview')}>
                Old Tasks
              </NavLink>
            </NavItem>
            <NavItem>
              <NavLink className={classnames({ active: activeTab === 'weeklyReview'})} onClick={()=> switchTab('weeklyReview')}>
                Weekly Review
              </NavLink>
            </NavItem>
            <NavItem>
              <NavLink className={classnames({ active: activeTab === 'monthlyReview'})} onClick={()=> switchTab('monthlyReview')}>
                Monthly Review
              </NavLink>
            </NavItem>
            <NavItem>
              <NavLink className={classnames({ active: activeTab === 'timeReport'})} onClick={()=> switchTab('timeReport')}>
                Time report
              </NavLink>
            </NavItem>
          </Nav>
          <Nav className="ml-auto" navbar>
            {(!auth.instance() || auth.instance().isAuthenticated())
              ? <NavItem onClick={() => { console.log('Logout'); auth.instance().logout() }}>
                <NavLink to="#" className="nav-link">
                  Logout</NavLink>
              </NavItem>
              : <NavItem onClick={() => { console.log('Login'); auth.instance().login() }}>
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
