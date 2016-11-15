/*
  input.js
  @flow

  This is the default <input type=text /> component.
*/

var React = require('react');

class Input extends React.Component {
  render() {
    return (
      <div style={styles.container}>
        <span style={styles.label}>{this.props.label}</span>
        <input type={this.props.type} style={styles.input} />
      </div>
    )
  }
}
Input.defaultProps = {
  type: 'text'
}

const styles = {
  container: {

  },
  label: {
    fontSize: 14,
    display: 'block'
  },
  input: {
    padding: 10,
    fontSize: 16,
    width: 300
  }
};

module.exports = Input;
