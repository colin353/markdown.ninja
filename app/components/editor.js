/*
  editor.js
  @flow

  The markdown editor, which uses ace.js as a code editor
  (see https://ace.c9.io for more details on that).
*/

var React = require('react');

type Props = {
  onChange: (markdown: string) => void
}

class Editor extends React.Component {
  editor: React.Component;
  ace: any;
  componentDidMount() {
    this.ace = window.ace.edit("editor");
    this.ace.setTheme("ace/theme/twilight");
    this.ace.session.setMode("ace/mode/markdown");
    this.ace.session.setUseWrapMode(true);
    
    // Register for onChange events.
    this.ace.on("change", this.textTyped.bind(this));
  }
  componentWillUnmount() {
    this.ace.destroy();
    this.ace.container.remove();
  }
  textTyped() {
    this.props.onChange(this.ace.getValue());
  }
  render() {
    return (
      <div style={styles.container}>
        <pre style={styles.editor}
          id="editor"
          ref={(e) => this.editor = e}
        ></pre>
      </div>
    );
  }
}

const styles = {
  container: {
    flex: 1
  },
  editor: {
    fontSize: 16,
    width: '100%',
    height: '100%',
    position:'relative',
    top:-16
  }
};

module.exports = Editor;
