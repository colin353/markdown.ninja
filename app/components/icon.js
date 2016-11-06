/*
  icon.js
  @flow

  A component for the Google Icon font.
*/

var React = require('react');

type Props = {
  name: string,
  style: { [key: string]: any }
};

class Icon extends React.Component {
  props: Props;
  static defaultProps: Props;

  render() {
    var style = Object.assign(styles.container, this.props.style);
    return (
      <i style={style} className="material-icons">{this.props.name}</i>
    )
  }
}

Icon.defaultProps = {
  name: "description",
  style: {}
};

const styles = {
  container: {
    position: 'relative',
    top: 5
  }
}

module.exports = Icon;
