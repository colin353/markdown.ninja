/*
  tree.js
  @flow

  Left-hand editor component, which acts like a file browser.
*/

var React = require('react');
import { ContextMenu, MenuItem, ContextMenuTrigger } from "react-contextmenu";

var Icon = require('./icon');
var Button = require('./button');
var Ellipsis = require('../tools/overflow-ellipsis');
var CSS = require('../tools/loadcss');
var Styles = require('../config/styles.json').styles;

import type { APIInstance } from '../api/api';

type Props = {
  api: APIInstance,
  pages: Page[],
  files: File[],
  setStyle: (s: string) => void,
  style: string,
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
      domain: this.props.api.user?this.props.api.user.domain:''
    };
  }

  componentDidMount() {
    this.props.api.addListener("authenticationStateChanged", "tree", () => {
      this.setState({ domain: this.props.api.user.domain });
    });
  }

  componentWillUnmount() {
    this.props.api.removeListeners("tree");
  }

  componentWillReceiveProps(props: Props) {
    if(props.style) {
        // Unload any existing styles.
        CSS.unload();
        // Add the current style to the page.
        CSS.load("/css/webstyles/" + props.style + ".css");
    }
  }

  handleClick() {

  }

  collect(page: Page) {
    return {page: page};
  }

  collectFile(file: File) {
    return {file: file};
  }

  selectStyle(name: string) {
    this.props.setStyle(name);
  }

  // Open your subdomain in a new tab.
  clickDomain() {
    window.open(window.location.protocol + "//" + this.state.domain + "." + this.props.api.BASE_DOMAIN, '_blank');
  }

  render() {
    console.log(Styles);
    return (
      <div style={styles.container}>
        <div onClick={this.clickDomain.bind(this)} style={styles.rootRow}><Icon name="book" /> {this.state.domain}.{this.props.api.BASE_DOMAIN}</div>
        {this.props.pages.map((p) => {
          return (
            <ContextMenuTrigger collect={this.collect.bind(this, p)} key={p.name} id="page">
              <div onClick={this.props.clickPage.bind(this, p)} className="noselect" style={styles.row}><Icon name="description" /> {p.name}</div>
            </ContextMenuTrigger>
          )
        })}

        {this.props.files.length?(
          <div style={styles.row}><Icon name="folder" /> files</div>
        ):[]}

        {this.props.files.map((f) => {
          return (
            <ContextMenuTrigger collect={this.collectFile.bind(this, f)} key={"page"+f.name} id="page">
              <div onClick={this.props.clickFile.bind(this, f)} className="noselect" style={styles.indentRow}><Icon name="description" /> {Ellipsis(f.name, 19)}</div>
            </ContextMenuTrigger>
          )
        })}

      <div>
        <div style={styles.row}><Icon name="format_paint" /> style
          <select value={this.props.style} onChange={(e) => this.selectStyle(e.target.value)} style={styles.select}>
            {Styles.map((s, index) => {
              return <option key={index}>{s}</option>;
            })}
          </select>
        </div>
      </div>

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
    paddingTop: 5,
    fontSize: 16,
    color: '#c4c4c4',
    backgroundColor: '#272822',
    width: 300,
    display: 'flex',
    flexDirection: 'column'
  },
  select: {
    marginLeft: 20,
    fontSize: 16,
    backgroundColor: '#61625E',
    color: 'rgb(196, 196, 196)',
    border: 'none',
    paddingLeft: 5,
    paddingRight: 5,
    paddingTop: 2,
    paddingBottom: 2
  },
  chooseStyle: {
    marginLeft: 60
  },
  rootRow: {
    paddingLeft: 20,
    cursor: 'default'
  },
  row: {
    paddingLeft: 40,
    cursor: 'default'
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
