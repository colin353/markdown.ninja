/*
  signup.js
  @flow

  A page for signing up.
*/

var React = require('react');

var Input = require('../components/input');
var Button = require('../components/button');

class Signup extends React.Component {
  render() {
    return (
      <div style={styles.container}>
        <h1>Sign up. </h1>
        <p>Already have an account? Sign in.</p>

        <Input label="name" />
        <Input label="domain" />
        <Input label="email" />
        <Input label="password" type="password" />

        <div style={{display: 'flex', marginTop: 10, marginBottom: 20, alignItems: 'center'}}>
          <a style={styles.forgot} href="#">forgot password?</a>
          <div style={{flex: 1}}></div>
          <Button color="red" action="log in" />
        </div>
      </div>
    );
  }
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
