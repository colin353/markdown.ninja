/*
  tree.js
  @flow

  Left-hand editor component, which acts like a file browser.
*/

var React = require('react');

import type { APIInstance } from '../api/api';
declare var api: APIInstance;

var Icon = require('./icon');

type Props = {
  [key: string]: any
}

class Tree extends React.Component {
  state: {
    domain: string
  };
  constructor(props: Props) {
    super(props);
    this.state = {
      domain: api.user?api.user.domain:''
    };
  }
  componentDidMount() {
    api.addListener("authenticationStateChanged", "tree", () => {
      this.setState({ domain: api.user.domain });
    });
  }
  componentWillUnmount() {
    api.removeListeners("tree");
  }
  render() {
    return (
      <div style={styles.container}>
        <div style={styles.rootRow}><Icon name="book" /> {this.state.domain}.{api.BASE_DOMAIN}</div>
        <div className="noselect" style={styles.row}><Icon name="description" /> index.md</div>
        <div className="noselect" style={styles.row}><Icon name="description" /> home.md</div>

      </div>
    );
  }
}

const styles = {
  container: {
    color: '#c4c4c4',
    backgroundColor: '#716669',
    width: 300,
    height: '100%'
  },
  rootRow: {
    paddingLeft: 20
  },
  row: {
    paddingLeft: 40
  }
};

module.exports = Tree;
