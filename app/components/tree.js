/*
  tree.js
  @flow

  Left-hand editor component, which acts like a file browser.
*/

var React = require('react');
import { ContextMenu, MenuItem, ContextMenuTrigger } from "react-contextmenu";

import type { APIInstance } from '../api/api';
declare var api: APIInstance;

var Icon = require('./icon');
var Button = require('./button');

type Props = {
  pages: Page[],
  files: File[],
  clickPage: (p: Page) => void,
  clickFile: (f: File) => void,
  onAddNewPage?: () => void,
  onUploadFile?: () => void
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
  handleClick() {

  }
  collect(page: Page) {
    return {page: page};
  }
  collectFile(file: File) {
    return {file: file};
  }
  render() {
    return (
      <div style={styles.container}>
        <div style={styles.rootRow}><Icon name="book" /> {this.state.domain}.{api.BASE_DOMAIN}</div>
        {this.props.pages.map((p) => {
          return (
            <ContextMenuTrigger collect={this.collect.bind(this, p)} key={p.name} id="page">
              <div onClick={this.props.clickPage.bind(this, p)} className="noselect" style={styles.row}><Icon name="description" /> {p.name}</div>
            </ContextMenuTrigger>
          )
        })}

        <div style={styles.row}><Icon name="folder" /> files</div>
        {this.props.files.map((f) => {
          return (
            <ContextMenuTrigger collect={this.collectFile.bind(this, f)} key={"page"+f.name} id="page">
              <div onClick={this.props.clickFile.bind(this, f)} className="noselect" style={styles.indentRow}><Icon name="description" /> {f.name}</div>
            </ContextMenuTrigger>
          )
        })}

      <div style={{flex: 1}}></div>
      <div style={styles.controlPanel}>
        <Button onClick={this.props.onAddNewPage} action="+ new page" />
        <div style={{marginLeft: 10}}></div>
        <Button onClick={this.props.onUploadFile} action="upload file" />
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
  indentRow: {
    paddingLeft: 60
  },
  controlPanel: {
    display: 'flex',
    marginBottom: 10,
    justifyContent: 'center'
  }
};

module.exports = Tree;
