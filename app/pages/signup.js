/*
  signup.js
  @flow

  A page for signing up.
*/

var React = require('react');

var Input = require('../components/input');
var Button = require('../components/button');

import type { APIInstance } from '../api/api';
declare var api: APIInstance;

class Signup extends React.Component {
  state: {
    name: string,
    domain: string,
    email: string,
    password: string,
    emailError: string,
    passwordError: string,
    domainError: string
  };

  constructor(props: {}) {
    super(props);

    this.state = {
      name: "",
      domain: "",
      email: "",
      password: "",
      emailError: '',
      passwordError: '',
      domainError: ''
    }
  }
  validate() {
    // Remove all errors.
    this.setState({
      emailError: '',
      passwordError: '',
      domainError: ''
    })

    // The domain is enforced to be valid using the setDomain function.
    // But you'll need to provide at least SOME domain.
    if(this.state.domain.length == 0) {
      this.setState({domainError: 'you must choose a domain'});
      return false;
    }

    // We do need to check the email: it's okay to leave it blank. But
    // if you do put an email address in, it should be valid.
    if(this.state.email != "" && !this.state.email.match(/^([a-zA-Z0-9_\-\.]+)@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.)|(([a-zA-Z0-9\-]+\.)+))([a-zA-Z]{2,4}|[0-9]{1,3})(\]?)$/)) {
      this.setState({emailError: 'that email address is invalid'})
      return false;
    }

    // Check that the password is long enough.
    if(this.state.password.length < 6) {
      this.setState({passwordError: 'your password is too short'})
      return false;
    }

    return true;
  }

  signUp() {
    // Try to validate the inputs.
    if(!this.validate()) return;

    // Now send the user information and signup.
    api.signup({
      email: this.state.email,
      password: this.state.password,
      name: this.state.name,
      domain: this.state.domain
    }).then(() => {
      // Signup worked.
      return api.checkAuth()
    }).then(() => {
      // Jump over to the "edit site" page.
      this.context.router.push("/edit/site");
    }).catch(() => {
      // Internal server error, I guess.
      this.setState({passwordError: "something went wrong creating account."})
    })
  }

  setDomain(domain: string) {
    // Ensure only valid domains are used.
    this.setState({domain: domain.toLowerCase().replace(/[^A-Za-z0-9_\\.]+/g, '')})
  }
  render() {
    return (
      <div style={styles.container}>
        <h1>Sign up. </h1>
        <p>Already have an account? <a href="#" onClick={this.context.router.push.bind(this.context.router, '/edit/login')}>Sign in.</a></p>

        <Input onReturn={this.signUp.bind(this)} value={this.state.name} onChange={(name) => this.setState({name})} label="name" />
        <Input onReturn={this.signUp.bind(this)} error={this.state.domainError} value={this.state.domain} onChange={this.setDomain.bind(this)} label="domain" />
        <Input onReturn={this.signUp.bind(this)} error={this.state.emailError} value={this.state.email} onChange={(email) => this.setState({email})} label="email" />
        <Input onReturn={this.signUp.bind(this)} error={this.state.passwordError} value={this.state.password} onChange={(password) => this.setState({password})} label="password" type="password" />

        <div style={{display: 'flex', marginTop: 10, marginBottom: 20, alignItems: 'center'}}>
          <a style={styles.forgot} href="#">forgot password?</a>
          <div style={{flex: 1}}></div>
          <Button onClick={this.signUp.bind(this)} color="red" action="sign up" />
        </div>
      </div>
    );
  }
}
Signup.contextTypes = {
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

module.exports = Signup;
