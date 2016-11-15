/*
  login.js
  @flow

  A page for logging in.
*/

var React = require('react');

var Input = require('../components/input');
var Button = require('../components/button');

class Login extends React.Component {
  render() {
    return (
      <div style={styles.container}>
        <h1>Login. </h1>

        <Input label="domain" />
        <Input label="password" type="password" />

        <div style={{display: 'flex', marginTop: 10, alignItems: 'center'}}>
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

module.exports = Login;
