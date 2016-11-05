/*
  titlebar.js
  @flow

  Renders the titlebar at the top of the application.
*/

var React = require('react');

import type { APIInstance } from '../api/api';
declare var api: APIInstance;

var Button = require('./button');
var Gravatar = require('./gravatar');

class Titlebar extends React.Component {
  state: {
    loggedIn: boolean,
    email: string
  };

  constructor(props: any) {
    super(props);

    this.state = {
      loggedIn: api.authenticated,
      email   : "test@fake_email_address.com"
    };
  }
  componentDidMount() {
    api.addListener("authenticationStateChanged", "titlebar", () => {
      console.log("authentication state changed.");
      this.setState({
        email   : api.user.email,
        loggedIn: api.authenticated
      });
    });
  }
  componentWillUnmount() {
    api.removeListeners("titlebar");
  }
  render() {
    return (
      <div style={styles.container}>
        <span style={styles.title}>Portfolio</span>
        <div style={styles.spacer}></div>
        {this.state.loggedIn?(
          <Gravatar email={this.state.email} />
        ):(
          <Button action="sign up" />
        )}
        <div style={{marginRight: 50}}></div>
      </div>
    );
  }
}

const styles = {
  container: {
    display: 'flex',
    flexDirection: 'row',
    height: 60,
    alignItems: 'center',
    borderBottom: '1px solid #CCC'
  },
  spacer: {
    flex: 1
  },
  title: {
    marginLeft: 50,
    fontWeight: 'bold',
    fontSize: '1.3em'
  }
}

module.exports = Titlebar;
