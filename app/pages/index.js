/*
  index.js
  @flow

  The homepage.
*/

var React = require('react');

var Button = require('../components/button');
var Icon = require('../components/icon');

class Index extends React.Component {
  clickSignup() { this.context.router.push('/edit/signup'); }
  render() {
    return (
      <div style={styles.container}>
        <h1>Build a website in 2 minutes flat.</h1>
        <p>
          Want a personal portfolio website? A simple blog? A website for any reason?
        </p>

        <p>Turn this:</p>

        <div style={styles.code}>
          <p># Hello world</p>
          <p>I'm a markdown ninja!</p>
          <p>![katana](files/katana.svg)</p>
        </div>

        <p>Into this:</p>

        <div style={{display: 'flex', flexDirection: 'row', alignItems: 'center', marginBottom: 20}}>
        <div style={styles.computer}>
          <div style={styles.addressBar}><Icon style={{fontSize: 16}} name="keyboard_backspace" /> <Icon style={{fontSize: 16}} name="close" /> <span style={{marginLeft: 10}}>http://yoursite.markdown.ninja/</span></div>
          <div style={styles.browser}>
            <h1>Hello world</h1>

            <p>I'm a markdown ninja!</p>

            <img height={50} src="/img/katana.svg" />
          </div>
        </div>

        <div className="phone" style={styles.phone}>
          <div style={styles.speaker}></div>
          <div style={styles.browser}>
            <h1>Hello world</h1>

            <p>I'm a markdown ninja!</p>

            <img height={50} src="/img/katana.svg" />
          </div>
          <div style={styles.homeButton}></div>
        </div>
        </div>

        <div style={styles.buttonRow}>
          <Button size="big" onClick={this.clickSignup.bind(this)} color="red" action="sign up free" />
        </div>
        <p style={styles.cred}>Thanks for <a style={{color: '#888'}} href="http://www.freepik.com/">Freepik</a> from <a style={{color: '#888'}} href="http://www.flaticon.com">Flatiron</a> for the icon.</p>

        <h1>Check out some examples</h1>
        <p>
          You can also type HTML elements directly into markdown, so you can make your site
          look pretty much however you want.
        </p>
        <br />

        <div style={styles.computerExample}>
          <div style={styles.addressBar}><Icon style={{fontSize: 16}} name="keyboard_backspace" /> <Icon style={{fontSize: 16}} name="close" /> <span style={{marginLeft: 10}}>http://colin.markdown.ninja/</span></div>
          <div style={styles.browserImage}>
            <img width={500} src="/img/example0.png" />
          </div>
        </div>

        <div style={styles.computerExample}>
          <div style={styles.addressBar}><Icon style={{fontSize: 16}} name="keyboard_backspace" /> <Icon style={{fontSize: 16}} name="close" /> <span style={{marginLeft: 10}}>http://jsmith.markdown.ninja/</span></div>
          <div style={styles.browserImage}>
            <img width={500} src="/img/example1.png" />
          </div>
        </div>

        <h1>Support this site</h1>

        <p>
          This website is 100% free for anyone to use. But it costs money to host the service.
        </p>

        <p>
          If you've benefitted from the service, <a href="https://www.paypal.me/markdownninja">consider
          donating to keep the service running</a>.
        </p>

        <p>
          If enough donations are received, I'll
          also buy a wildcard SSL certificate, so that all of the subdomains can support SSL.
        </p>

      </div>
    );
  }
}
Index.contextTypes = {
  'router': React.PropTypes.object
};

const styles = {
  container: {
    maxWidth: 700,
    marginLeft: 'auto',
    marginRight: 'auto',
    display: 'block'
  },
  code: {
    backgroundColor: 'white',
    padding: 10,
    paddingLeft: 20,
    borderRadius: 5,
    paddingRight: 20,
    marginLeft: 'auto',
    maxWidth: 400,
    marginRight: 'auto',
    boxShadow: 'inset 0px 0px 5px #AAA',
    marginBottom: 30
  },
  computer: {
    backgroundColor: '#F0F0F0',
    border: '2px solid #333',
    borderTopLeftRadius: 5,
    borderTopRightRadius: 5,
    marginLeft: 'auto',
    minWidth: 350,
    marginRight: 'auto'
  },
  computerExample: {
    display: 'block',
    width: 500,
    backgroundColor: '#F0F0F0',
    border: '2px solid #333',
    borderTopLeftRadius: 5,
    borderTopRightRadius: 5,
    marginLeft: 'auto',
    marginRight: 'auto',
    marginBottom: 40
  },
  addressBar: {
    padding: 10,
    paddingTop: 5,
    fontSize: 12
  },
  browser: {
    borderTop: '2px solid #333',
    backgroundColor: 'white',
    padding: 10,
    paddingLeft: 20,
    paddingRight: 20
  },
  browserImage: {
    borderTop: '2px solid #333',
    backgroundColor: 'white'
  },
  buttonRow: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center'
  },
  phone: {
    width: 165,
    marginLeft: 'auto',
    marginRight: 'auto',
    paddingLeft: 4,
    paddingRight: 5,
    backgroundColor: '#333',
    paddingTop: 15,
    paddingBottom: 10,
    borderRadius: 12
  },
  speaker: {
    marginLeft: 'auto',
    marginRight: 'auto',
    backgroundColor: '#222',
    width: 50,
    height: 7,
    borderRadius: 15,
    marginBottom: 15
  },
  homeButton: {
    marginLeft: 'auto',
    marginRight: 'auto',
    backgroundColor: '#222',
    width: 30,
    height: 30,
    marginTop: 10,
    borderRadius: 50
  },
  cred: {
    marginTop: 40,
    textAlign: 'center',
    fontSize: 12,
    color: '#888'
  }
}

module.exports = Index;
