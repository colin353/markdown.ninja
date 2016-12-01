/*
  tree.js
  @flow

  Left-hand editor component, which acts like a file browser.
*/

var React = require('react');
var ReactDOM = require('react-dom');
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
  onUploadFile?: () => void,
  micro: boolean
}

class Tree extends React.Component {
  props: Props;
  state: {
    domain: string,
    collapsed: boolean
  };
  constructor(props: Props) {
    super(props);
    this.state = {
      domain: this.props.api.user?this.props.api.user.domain:'',
      collapsed: true
    };
  }

  componentDidMount() {
    this.props.api.addListener("authenticationStateChanged", "tree", () => {
      this.setState({ domain: this.props.api.user.domain });
    });

    this.props.api.addListener('clickBody', 'tree', (e) => {
      var area = ReactDOM.findDOMNode(this);
      if (!area.contains(e.target)) {
        if(!this.state.collapsed) this.setState({collapsed: true});
      }
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
    if(this.props.micro && this.state.collapsed) {
      return (<div onClick={() => this.setState({collapsed: false})} className="noselect" style={styles.collapsedContainer}><Icon style={styles.menuIcon} name="menu" /></div>);
    }
    return (
      <div style={this.props.micro?styles.microcontainer:styles.container}>
        <div style={{flex: 1, overflow: 'auto'}}>
        <div onClick={this.clickDomain.bind(this)} style={styles.rootRow}><Icon name="book" /> {Ellipsis(this.state.domain+"."+this.props.api.BASE_DOMAIN, 25)}</div>
        {this.props.pages.map((p) => {
          return (
            <ContextMenuTrigger collect={this.collect.bind(this, p)} key={p.name} id="page">
              <div onClick={this.props.clickPage.bind(this, p)} className="noselect" style={styles.row}><Icon name="description" /> {Ellipsis(p.name, 19)}</div>
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
      </div>
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
  collapsedContainer: {
    paddingTop: 5,
    fontSize: 16,
    color: '#c4c4c4',
    backgroundColor: '#272822',
    width: 50,
    display: 'flex',
    flexDirection: 'column'
  },
  microcontainer: {
    paddingTop: 5,
    fontSize: 16,
    color: '#c4c4c4',
    backgroundColor: '#272822',
    width: 300,
    marginRight: -250,
    zIndex: 5,
    height: '100%',
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
    marginTop: 5,
    marginBottom: 5,
    paddingLeft: 20,
    cursor: 'default'
  },
  row: {
    paddingLeft: 40,
    cursor: 'default',
    marginTop: 5,
    marginBottom: 5
  },
  indentRow: {
    paddingLeft: 60,
    marginTop: 5,
    marginBottom: 5
  },
  controlPanel: {
    display: 'flex',
    marginBottom: 10,
    justifyContent: 'center'
  },
  menuIcon: {
    marginLeft: 13,
    marginTop: 5
  }
};

module.exports = Tree;
