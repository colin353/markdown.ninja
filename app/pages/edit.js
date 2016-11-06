/*
  edit.js
  @flow

  The site editor page.
*/

var React = require('react');

import type { APIInstance } from '../api/api';
declare var api: APIInstance;

var Tree = require('../components/tree');
var Editor = require('../components/editor');
var Preview = require('../components/preview');

class Edit extends React.Component {
  state: {
    markdown: string,
    html: string
  };

  converter: any;

  constructor(props: any) {
    super(props);

    this.state = {
      html: "",
      markdown: ""
    }

    this.converter = new window.showdown.Converter();
  }

  textEdited(markdown: string) {
    var html = this.converter.makeHtml(markdown);
    this.setState({markdown, html});
  }
  render() {
    return (
      <div style={styles.container}>
        <Tree />
        <Editor onChange={this.textEdited.bind(this)} />
        <Preview html={this.state.html} />
      </div>
    );
  }
}

const styles = {
  container: {
    display: 'flex',
    flexDirection: 'row',
    flex: 1
  }
};

module.exports = Edit;
