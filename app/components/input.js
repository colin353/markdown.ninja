/*
  input.js
  @flow

  This is the default <input type=text /> component.
*/

var React = require('react');

type Props = {
  label: string,
  onChange: (val: string) => void,
  onReturn: () => void,
  type?: string,
  value: string,
  error: string,
  success: string,
};

class Input extends React.Component {
  props: Props;
  static defaultProps: Props;

  onChange(e: {target: {value: string}}) {
    this.props.onChange(e.target.value);
  }
  onKeyDown(e: any) {
    if(e.key == "Enter") this.props.onReturn();
  }
  render() {
    return (
      <div style={styles.container}>
        <span style={styles.label}>{this.props.label}</span>
        <input onKeyDown={this.onKeyDown.bind(this)} value={this.props.value} type={this.props.type} onChange={this.onChange.bind(this)} style={styles.input} />
        {this.props.error?(
          <div style={styles.error}>{this.props.error}</div>
        ):this.props.success?(
          <div style={styles.success}>{this.props.success}</div>
        ):[]}
      </div>
    )
  }
}
Input.defaultProps = {
  value: '',
  label: '',
  type: 'text',
  onChange: () => {},
  onReturn: () => {},
  error: '',
  success: ''
}

const styles = {
  container: {
  },
  error: {
    backgroundColor: 'rgb(201, 137, 134)',
    color: 'white',
    padding: 5,
    marginLeft: 1,
    position: 'relative',
    top: -1,
    width: 312,
    fontSize: 14,
    textAlign: 'center',
    borderBottomLeftRadius: 3,
    borderBottomRightRadius: 3
  },
  success: {
    backgroundColor: '#7A89C2',
    color: 'white',
    padding: 5,
    marginLeft: 1,
    position: 'relative',
    top: -1,
    width: 312,
    fontSize: 14,
    textAlign: 'center',
    borderBottomLeftRadius: 3,
    borderBottomRightRadius: 3
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
