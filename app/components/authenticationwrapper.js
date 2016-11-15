/*
  authenticationwrapper.js
  @flow

  This component wraps pages which require authentication. If the user ever
  loses authentication, it'll automatically redirect them back to the login
  page.
*/

var React = require('react');

import type { APIInstance } from '../api/api';
declare var api: APIInstance;

class AuthenticationWrapper extends React.Component {
  componentDidMount() {
    // If the authentication state changes, and we lose authentication,
    // redirect to the login page.
    api.addListener("authenticationStateChanged", "authwrapper", (newState) => {
      if(!newState) this.context.router.push('/edit/login');
    })

    // Also, check if we are not initially authenticated.
    if(!api.authenticated) this.context.router.push('/edit/login');
  }
  componentWillUnmount() {
    api.removeListeners("authwrapper");
  }
  render() {
    return this.props.children;
  }
}
AuthenticationWrapper.contextTypes = {
  router: React.PropTypes.object
}

module.exports = AuthenticationWrapper;
