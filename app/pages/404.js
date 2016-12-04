/*
  404.js
  @flow

  A page for 404s.
*/

var React = require('react');

class YouAreLost extends React.Component {
  render() {
    return (
      <div style={styles.container}>
        <h1>404.</h1>
        <p>That page doesn't exist. Did you get here by accident?</p>
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

module.exports = YouAreLost;
