/*
  404.js
  @flow

  A page for 404s.
*/

var React = require('react');

class YouAreLost extends React.Component {
  render() {
    return (
      <div>
        <h1>404.</h1>
        <p>That page doesn't exist. Did you get here by accident?</p>
      </div>
    );
  }
}

module.exports = YouAreLost;
