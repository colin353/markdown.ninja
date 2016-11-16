/*
  index.js
  @flow

  The homepage.
*/

var React = require('react');

var Button = require('../components/button');

class Index extends React.Component {
  render() {
    return (
      <div style={styles.container}>
        <h1>Build a personal website in 2 minutes flat.</h1>
        <p>
          These days, a resume isn't enough. You have to show off what you've done.
          But maintaining a portfolio website is time consuming and annoying.
        </p>

        <p>
          With portfolio, you can make a personal website in 2 minutes with the magic
          of markdown.
        </p>

        <div style={styles.buttonRow}>
          <Button onClick={this.context.router.push.bind(this.context.router, '/edit/signup')} color="red" action="sign up free" />
        </div>
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
  buttonRow: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center'
  }
}

module.exports = Index;
