/*
  login.js
  @flow

  A page for logging in.
*/

var React = require('react');

var Input = require('../components/input');
var Button = require('../components/button');

import type { APIInstance } from '../api/api';
declare var api: APIInstance;

class Login extends React.Component {
  state: {
    domain: string,
    password: string,
    passwordError: string
  };

  constructor(props: {}) {
    super(props);

    this.state = {
      domain: '',
      password: '',
      passwordError: ''
    }
  }

  login() {
    // Clear any existing errors.
    this.setState({passwordError: ''});

    // Try to log in.
    api.login(this.state.domain, this.state.password).then(() => {
      // Login must be successful if we reached this point.
      // Now we'll route to the edit site page.
      this.context.router.push('/edit/site');
    }).catch(() => {
      // An error occurred, report it.
      this.setState({passwordError: "that domain/password didn't work" });
    })
  }

  render() {
    return (
      <div style={styles.container}>
        <h1>Login. </h1>
        <p>Don't have an account? <a href="#" onClick={this.context.router.push.bind(this.context.router, '/edit/signup')}>Sign up.</a></p>

        <Input onReturn={this.login.bind(this)} value={this.state.domain} onChange={(domain) => this.setState({domain})} label="domain" />
        <Input onReturn={this.login.bind(this)} error={this.state.passwordError} value={this.state.password} onChange={(password) => this.setState({password})} label="password" type="password" />

        <div style={{display: 'flex', marginTop: 10, alignItems: 'center'}}>
          <a style={styles.forgot} href="#">forgot password?</a>
          <div style={{flex: 1}}></div>
          <Button onClick={this.login.bind(this)} color="red" action="log in" />
        </div>
      </div>
    );
  }
}
Login.contextTypes = {
  router: React.PropTypes.object
}

const styles = {
  container: {
    width: 320,
    marginLeft: 'auto',
    marginRight: 'auto',
    display: 'block'
  },
  forgot: {
    fontSize: 14,
    textDecoration: 'none',
    color: '#aaa'
  }
}

module.exports = Login;
