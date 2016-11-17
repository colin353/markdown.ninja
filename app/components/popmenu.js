/*
  popmenu.js
  @flow

  A menu which pops up and displays user information,
  lets you log out, etc.
*/

var React = require('react');

var Button = require('./button');

import type { APIInstance } from '../api/api';

type Props = {
  api: APIInstance
};

class PopMenu extends React.Component {
  props: Props;
  state: {
    domain: string
  };

  constructor(props: Props) {
    super(props);
    this.state = {
      domain: this.props.api.user?this.props.api.user.domain:""
    };
  }
  componentDidMount() {
    this.props.api.addListener("authenticationStateChanged", "popmenu", () => {
      this.setState({
        domain: this.props.api.user?this.props.api.user.domain:""
      });
    });
  }
  componentWillUnmount() {
    this.props.api.removeListeners("popmenu");
  }

  logout() {
    this.props.api.logout().then(() => {
      this.context.router.push('/edit/login');
    });
  }

  clickEditSite() {
    this.context.router.push('/edit/site');
  }

  // When the user clicks on their domain, we'll open it for them
  // in a new tab.
  clickDomain() {
    window.open(window.location.protocol + "//" + this.state.domain + "." + this.props.api.BASE_DOMAIN, '_blank');
  }

  render() {
    return (
      <div style={styles.container}>
        <div style={styles.arrow}></div>
        <div onClick={this.clickDomain.bind(this)} style={styles.domain}>
          <span style={styles.subdomain}>{this.state.domain}</span>
          <span>{"." + this.props.api.BASE_DOMAIN}</span>
        </div>
        <div style={styles.buttonlist}>
          <Button color='red' onClick={this.clickEditSite.bind(this)} action="edit my site" />
          <div style={{display: 'flex', flex: 1}}></div>
          <Button onClick={this.logout.bind(this)} action="log out" />
        </div>
      </div>
    );
  }
}

PopMenu.contextTypes = {
  router: React.PropTypes.object
}

const styles = {
  buttonlist: {
    display: 'flex',
    flexDirection: 'row',
    margin: 10
  },
  container: {
    position: 'absolute',
    right: -20,
    top: 10,
    width: 300,
    backgroundColor: '#c4c4c4',
    zIndex: 10
  },
  domain: {
    cursor: 'default',
    flex: 1,
    display: 'flex',
    alignItems: 'center',
    paddingTop: 20,
    paddingBottom: 20,
    justifyContent: 'center',
    boxShadow: '#AAA 0px 0px 5px inset',
    marginLeft: 10,
    marginRight: 10,
    marginTop: 10
  },
  subdomain: {
    fontWeight: 'bold'
  },
  arrow: {
    width: 0,
    height: 0,
    borderLeft: '10px solid transparent',
    borderRight: '10px solid transparent',
    borderBottom: '10px solid #c4c4c4',
    marginTop:-8,
    marginLeft:250
  }
}

module.exports = PopMenu;
