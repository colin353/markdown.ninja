/*
  editor.js
  @flow

  The markdown editor, which uses ace.js as a code editor
  (see https://ace.c9.io for more details on that).
*/

var React = require('react');

type Props = {
  onChange: (markdown: string) => void,
  visible: boolean,
  initialText: string
}

class Editor extends React.Component {
  editor: React.Component<any, any, any>;
  ace: any;

  componentDidMount() {
    this.ace = window.ace.edit("editor");
    this.ace.setTheme("ace/theme/custom");
    this.ace.session.setMode("ace/mode/markdown");
    this.ace.session.setUseWrapMode(true);

    // Not sure what this does, but I was asked to do it by the
    // console window.
    this.ace.$blockScrolling = Infinity;

    // Register for onChange events.
    this.ace.on("change", this.textTyped.bind(this));

    // Try to put the initial text on there.
    if(this.props.initialText) this.ace.session.setValue(this.props.initialText, -1);
  }
  componentWillUnmount() {
    this.ace.destroy();
    this.ace.container.remove();
  }
  textTyped() {
    this.props.onChange(this.ace.getValue());
  }
  componentWillReceiveProps(nextProps: Props) {
    if(nextProps.initialText != this.props.initialText) {
      this.ace.setValue(nextProps.initialText, -1);
    }
  }
  render() {
    var containerstyle = Object.assign({}, styles.container);
    if(!this.props.visible) containerstyle.display = 'none';
    return (
      <div style={containerstyle}>
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
