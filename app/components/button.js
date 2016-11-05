/*
  button.js
  @flow

  The standard button.
*/

var React = require('react');

class Button extends React.Component {
  render() {
    return (
      <div className="button" style={styles.container}>
        <span style={styles.text}>{this.props.action}</span>
      </div>
    );
  }
}

const styles = {
  container: {
    backgroundColor: '#7A89C2',
    paddingLeft: 10,
    paddingRight: 10,
    paddingTop: 5,
    paddingBottom: 5,
    borderRadius: 2
  },
  text: {
    color: '#E0E0E0'
  }
}

module.exports = Button;
