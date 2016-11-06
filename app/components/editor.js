/*
  editor.js
  @flow

  The markdown editor, which uses ace.js as a code editor
  (see https://ace.c9.io for more details on that).
*/

var React = require('react');

class Editor extends React.Component {
  editor: React.Component;
  ace: any;
  componentDidMount() {
    this.ace = window.ace.edit("editor");
    this.ace.setTheme("ace/theme/twilight");
    this.ace.session.setMode("ace/mode/markdown");
  }
  render() {
    return (
      <div style={styles.container}>
        <pre style={styles.editor} id="editor" ref={(e) => this.editor = e}>This is an editor.</pre>
      </div>
    );
  }
}

const styles = {
  container: {
    flex: 1
  },
  editor: {
    width: '100%',
    height: '100%',
    position:'relative',
    top:-12
  }
};

module.exports = Editor;
