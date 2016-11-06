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

class Edit extends React.Component {
  render() {
    return (
      <div style={styles.container}>
        <Tree />
        <Editor />
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
