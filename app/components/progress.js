/*
  progress.js

  @flow
*/

var React = require('react');

type Props = {
  progress: number
}

class Progress extends React.Component {
  render() {
    var indicator = Object.assign({}, { width: this.props.progress+"%" }, styles.indicator);
    console.log(indicator);
    return (
      <div style={styles.container}>
        <div style={indicator}></div>
        <div style={styles.text}>{this.props.progress.toFixed(0)}%</div>
      </div>
    );
  }
}

const styles = {
  indicator: {
    display: 'block',
    height: 50,
    backgroundColor: 'rgb(73, 72, 62)'
  },
  container: {
    marginTop: 20,
    backgroundColor: 'rgb(51, 51, 51)'
  },
  text: {
    color: 'white',
    textAlign: 'center',
    marginTop: -20,
    position: 'relative',
    top: -13,
  }
}

module.exports = Progress;
