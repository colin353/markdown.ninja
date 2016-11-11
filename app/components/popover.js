/*
  popover.js
  @flow

  This component renders a popover, which, when visible, dims the rest of the screen
  and appears in the center, like a dialog box.
*/

var React = require('react');

import type { APIInstance } from '../api/api';
declare var api: APIInstance;

type Props = {
  visible: boolean,
  onDismiss: Function,
  onFocus: Function,
  children?: Array<React.Component>
}

class Popover extends React.Component {
  props: Props;
  static defaultProps: Props;

  componentDidMount() {
    api.addListener("escapeKeyPressed", "popover", () => {
      this.props.onDismiss();
    })
  }
  componentWillUnmount() {
    api.removeListeners("popover");
  }
  componentWillReceiveProps(newprops: Props) {
    if(newprops.visible && !this.props.visible) {
      setTimeout(this.props.onFocus, 300);
    }
  }
  render() {
    var shade = Object.assign({}, styles.shade, this.props.visible?{display: 'block'}:{});
    var popover = Object.assign({}, styles.popover, this.props.visible?{display: 'block'}:{});
    return (
      <div>
        <div onClick={this.props.onDismiss} style={shade}></div>
        <div style={popover}>{this.props.children}</div>
      </div>
    );
  }
}
Popover.defaultProps = {
  visible: false,
  onDismiss: () => {},
  onFocus: () => {}
}

const styles = {
  shade: {
    display: 'none',
    backgroundColor: '#333',
    opacity: 0.2,
    position: 'absolute',
    left: 0,
    top: 0,
    width: '100%',
    height: '100%',
    zIndex: 20
  },
  popover: {
    display: 'none',
    padding: 20,
    backgroundColor: '#E0E0E0',
    width: 400,
    height: 100,
    position: 'absolute',
    zIndex: 25,
    top: '20%',
    left: '50%',
    marginLeft: -200
  }
};

module.exports = Popover;
