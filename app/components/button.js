/*
  button.js
  @flow

  The standard button.
*/

var React = require('react');

type Props = {
  action: string,
  color: 'red' | 'blue',
  size: 'big' | 'small' | 'wide',
  onClick: Function
};

class Button extends React.Component {
  props: Props;
  static defaultProps: Props;

  render() {
    var buttonstyle = Object.assign({}, types[this.props.size]);
    buttonstyle.backgroundColor = colors[this.props.color];

    return (
      <div onClick={this.props.onClick} className="button" style={buttonstyle}>
        <span style={styles.text}>{this.props.action}</span>
      </div>
    );
  }
}

Button.propTypes = {
  action: React.PropTypes.string.isRequired,
  onClick: React.PropTypes.func
};

Button.defaultProps = {
  action: '',
  color: 'blue',
  onClick: () => {},
  size: 'small'
};

const colors = {
  'blue': '#7A89C2',
  'red': '#C98986'
}

const styles = {
  container: {
    textAlign: 'center',
    display: 'inline-block',
    backgroundColor: '#7A89C2',
    paddingLeft: 10,
    paddingRight: 10,
    paddingTop: 5,
    paddingBottom: 5,
    borderRadius: 2
  },
  bigContainer: {
    textAlign: 'center',
    display: 'inline-block',
    backgroundColor: '#7A89C2',
    paddingLeft: 20,
    paddingRight: 20,
    paddingTop: 10,
    paddingBottom: 10,
    fontSize: 18,
    borderRadius: 2
  },
  wideContainer: {
    textAlign: 'center',
    display: 'block',
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

const types = {
  big: styles.bigContainer,
  small: styles.container,
  wide: styles.wideContainer
}

module.exports = Button;
