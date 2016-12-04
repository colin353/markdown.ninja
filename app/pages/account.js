/*
  account.js
  @flow

  This page lets you manage your account, doing things like resetting your
  password, deleting the account, registering a domain, etc.
*/

var React = require('react');

var Input = require('../components/input');
var Button = require('../components/button');

import type { APIInstance } from '../api/api';

type Props = {
  api: APIInstance
}

class AccountPage extends React.Component {
  state: {
    domain: string,
    email: string,
    password: string,
    password2: string,
    emailError: string,
    emailSuccess: string,
    customDomain: string,
    customDomainSuccess: string,
    customDomainError: string
  }

  constructor(props: Props) {
    super(props);
    this.state = {
      domain: this.props.api.user?this.props.api.user.domain:"",
      email: this.props.api.user?this.props.api.user.email:"",
      emailSuccess: "",
      emailError: "",
      password: "",
      password2: "",
      passwordError: '',
      passwordSuccess: '',
      customDomain: this.props.api.user?this.props.api.user.external_domain:"",
      customDomainError: '',
      customDomainSuccess: ''
    };
  }

  validateEmail() {
    // We do need to check the email: it's okay to leave it blank. But
    // if you do put an email address in, it should be valid.
    if(this.state.email != "" && !this.state.email.match(/^([a-zA-Z0-9_\-\.]+)@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.)|(([a-zA-Z0-9\-]+\.)+))([a-zA-Z]{2,4}|[0-9]{1,3})(\]?)$/)) {
      this.setState({emailError: 'that email address is invalid'})
      return false;
    }

    this.setState({emailError: ''});

    return true;
  }

  validatePassword() {
    // Check that the password is long enough.
    if(this.state.password.length < 6) {
      this.setState({passwordError: 'your password is too short'});
      return false;
    }
    if(this.state.password != this.state.password2) {
      this.setState({passwordError: "your passwords don't match"});
      return false;
    }

    this.setState({passwordError: ''});

    return true;
  }

  componentDidMount() {
    this.props.api.addListener("authenticationStateChanged", "popmenu", () => {
      this.setState({
        domain: this.props.api.user?this.props.api.user.domain:"",
        email: this.props.api.user?this.props.api.user.email:"",
        customDomain: this.props.api.user?this.props.api.user.external_domain:""
      });
    });
  }

  updateEmail() {
    if(!this.validateEmail()) return;

    this.props.api.updateEmail(this.state.email).then((r) => {
      if(r.error) this.setState({emailError: 'that email address is invalid'});
      else this.setState({emailSuccess: 'email address updated'});
    })
  }

  updatePassword() {
    if(!this.validatePassword()) return;

    this.props.api.updatePassword(this.state.password).then((r) => {
      if(r.error) this.setState({passwordError: "couldn't update password"});
      else this.setState({passwordSuccess: 'password updated'});
    })
  }

  validateDomain() {
    return true;
  }

  updateCustomDomain() {
    this.setState({
      customDomainError: '',
      customDomainSuccess: ''
    });

    this.props.api.updateCustomDomain(this.state.customDomain).then((r) => {
      if(r.error) {
        if(r.result == "duplicate") this.setState({customDomainError: "that domain is taken"});
        else this.setState({customDomainError: "that domain is invalid"});
      }
      else this.setState({customDomainSuccess: 'custom domain set'});
    }).catch((r) => {
      this.setState({customDomainError: "that domain is invalid or taken"});
    })
  }

  render() {
    return (
      <div style={styles.container}>
        <h3>Your account details</h3>
        <p>Your website is: <a href={"http://" + this.state.domain + "." + this.props.api.BASE_DOMAIN}>{"http://" + this.state.domain + "." + this.props.api.BASE_DOMAIN}</a></p>
        <Input success={this.state.emailSuccess} value={this.state.email} onChange={(e) => this.setState({email: e})} label="email" error={this.state.emailError} />
        <p>You don't need to set an email address. It's only used in case you forget your password.</p>
        <Button onClick={this.updateEmail.bind(this)} action="update email" />
        <h3>Set a new password</h3>
        <Input value={this.state.password} onChange={(p) => this.setState({password:p})} type="password" label="password" />
        <Input error={this.state.passwordError} success={this.state.passwordSuccess} value={this.state.password2} onChange={(p) => this.setState({password2: p})} type="password" label="retype password" />
        <p></p>
        <Button action="update password" onClick={this.updatePassword.bind(this)} />
        <h3>Use a domain you own</h3>
        <p>
          It's easy to set up a domain you already own to point to your markdown.ninja site. Just type the domain you own
          into the box below, then set up a CNAME entry in your domain's DNS settings for markdown.ninja.
        </p>
        <Input value={this.state.customDomain} error={this.state.customDomainError} success={this.state.customDomainSuccess} onChange={(d) => this.setState({customDomain: d})} label="your custom domain" />
        <p></p>
        <Button onClick={this.updateCustomDomain.bind(this)} action="set custom domain" />
      </div>
    );
  }
}

const styles = {
  container: {
    maxWidth: 700,
    marginLeft: 'auto',
    marginRight: 'auto'
  }
}

module.exports = AccountPage;
