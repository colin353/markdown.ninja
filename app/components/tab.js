/*
  tab.js
  @flow

  A top tab bar for the editor.
*/

var React = require('react');
var Ellipsis = require('../tools/overflow-ellipsis');

type Props = {
  selected: boolean,
  name: string,
  onClick?: Function,
  indicator?: boolean
}

class Tab extends React.Component {
  props: Props;
  static defaultProps: Props;

  render() {
    var containerstyle = Object.assign({}, styles.container, this.props.selected?styles.selected:{});
    return (
      <div onClick={this.props.onClick} className="noselect" style={containerstyle}>
        <span>
          {Ellipsis(this.props.name, 16)}
          {this.props.indicator?(
            <div style={styles.indicator}></div>
          ):[]}
        </span>
      </div>
    );
  }
}
Tab.defaultProps = {
  name: 'untitled.md',
  selected: false,
  indicator: false
}

const styles = {
  container: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: 5,
    width: 200,
    height: 40,
    color: 'rgb(196, 196, 196)',
    backgroundColor: '#272822',
    borderTopRightRadius: 3,
    borderTopLeftRadius: 3,
    borderLeft: '2px solid #222',
    borderTop: '2px solid #222',
    borderRight: '2px solid #222',
    position: 'relative',
    top: 2,
    zIndex: 3
  },
  indicator: {
    backgroundColor: '#7A89C2',
    height: 5,
    width: 5,
    color: 'white',
    display: 'inline-block',
    marginLeft: 10,
    borderRadius: 5,
    position: 'relative',
    top: -2
  },
  selected: {
    backgroundColor: '#49483E'
  }
}

module.exports = Tab;
