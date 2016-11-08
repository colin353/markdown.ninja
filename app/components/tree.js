/*
  tree.js
  @flow

  Left-hand editor component, which acts like a file browser.
*/

var React = require('react');

import type { APIInstance } from '../api/api';
declare var api: APIInstance;

var Icon = require('./icon');
var Button = require('./button');

type Props = {
  pages: Page[],
  clickPage: (p: Page) => void
}

class Tree extends React.Component {
  props: Props;
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
        {this.props.pages.map((p) => {
          return (
            <div onClick={this.props.clickPage.bind(this, p)} key={p.name} className="noselect" style={styles.row}><Icon name="description" /> {p.name}</div>
          )
        })}

      <div style={{flex: 1}}></div>
      <div style={styles.controlPanel}>
        <Button action="+ new page" />
        <div style={{marginLeft: 10}}></div>
        <Button action="upload file" />
      </div>
      </div>
    );
  }
}

const styles = {
  container: {
    fontSize: 16,
    color: '#c4c4c4',
    backgroundColor: '#716669',
    width: 300,
    height: '100%',
    display: 'flex',
    flexDirection: 'column'
  },
  rootRow: {
    paddingLeft: 20
  },
  row: {
    paddingLeft: 40
  },
  controlPanel: {
    display: 'flex',
    marginBottom: 10,
    justifyContent: 'center'
  }
};

module.exports = Tree;
