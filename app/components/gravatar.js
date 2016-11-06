/*
  gravatar.js
  @flow

  Display a gravatar image based upon a user's
  email address.
*/

var React = require('react');

type Props = {
  email: string,
  onClick: Function
}

class Gravatar extends React.Component {
  props: Props;
  static defaultProps: Props;

  render() {
    var hash = window.md5(this.props.email);

    return (
      <div onClick={this.props.onClick} style={styles.container}>
        <img style={styles.image} src={"https://www.gravatar.com/avatar/"+hash+"?d=mm&s=80"} />
      </div>
    );
  }
}
Gravatar.propTypes = {
  email: React.PropTypes.string,
  onClick: () => {}
}

const styles = {
  container: {
    backgroundColor: '#c4c4c4',
    boxShadow: '2px 1px 1px #AAA',
    borderRadius: 100,
    width: 40,
    height: 40,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center'
  },
  image: {
    width: 35,
    height: 35,
    borderRadius: 100
  }
};

module.exports = Gravatar;
