/*
  titlebar.js
  @flow

  Renders the titlebar at the top of the application.
*/

var React = require('react');
var ReactDOM = require('react-dom');

var Button = require('./button');
var Gravatar = require('./gravatar');
var PopMenu = require('./popmenu');

import type { APIInstance } from '../api/api';

type Props = {
  api: APIInstance
};

class Titlebar extends React.Component {
  state: {
    loggedIn: boolean,
    email: string,
    menuVisible: boolean
  };

  menu: PopMenu;
  gravatar: Gravatar;

  constructor(props: any) {
    super(props);

    this.state = {
      loggedIn   : this.props.api.authenticated,
      email      : "test@fake_email_address.com",
      menuVisible: false
    };
  }
  componentDidMount() {
    this.props.api.addListener("authenticationStateChanged", "titlebar", () => {
      this.setState({
        email   : this.props.api.user?this.props.api.user.email:"",
        loggedIn: this.props.api.authenticated
      });
    });

    // This is a bit of a hack to make the menu work the way you would expect.
    // When the user clicks on the body of the page, but NOT inside of the
    // menu, we should close the menu. If they click on the gravatar, we should
    // let the gravatar handle closing the menu (otherwise you'll get a close -> open)
    // as both handlers run. So we must detect the area of the click and only close
    // the menu if we are clicking away from both the menu itself and the gravatar.
    this.props.api.addListener("clickBody", "titlebar", (e) => {
      if(!this.gravatar) return;

      var isClickingMenu = false;
      if(this.menu) {
        var menuArea = ReactDOM.findDOMNode(this.menu);
        isClickingMenu = menuArea.contains(e.target);
      }

      var gravatarArea = ReactDOM.findDOMNode(this.gravatar);
      if (!isClickingMenu && !gravatarArea.contains(e.target)) {
        this.setState({ menuVisible: false });
      }
    });
  }
  componentWillUnmount() {
    this.props.api.removeListeners("titlebar");
  }
  toggleMenu() {
    this.setState({
      menuVisible: !this.state.menuVisible
    });
  }
  render() {
    return (
      <div style={styles.container}>
        <span onClick={this.context.router.push.bind(this.context.router, '/')} style={styles.title} className="noselect">Portfolio</span>
        <div style={styles.spacer}></div>
        {this.state.loggedIn?(
          <div>
            <Gravatar ref={(g) => this.gravatar = g} onClick={this.toggleMenu.bind(this)} email={this.state.email} />
            <div style={{position: 'relative'}}>
              {this.state.menuVisible?(
                  <PopMenu api={this.props.api} ref={(m) => this.menu = m}  />
              ):[]}
            </div>
          </div>
        ):(
          <Button onClick={this.context.router.push.bind(this, '/edit/signup')} action="sign up" />
        )}
        <div style={{marginRight: 50}}></div>
      </div>
    );
  }
}

Titlebar.contextTypes = {
  router: React.PropTypes.object
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
