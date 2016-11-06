/*
  preview.js
  @flow

  Renders the HTML output of the markdown.
*/

var React = require('react');

type Props = {
  html: string
}

class Preview extends React.Component {
  render() {
    return (
      <div style={styles.container} dangerouslySetInnerHTML={{ __html: this.props.html}}></div>
    );
  }
}

const styles = {
  container: {
    flex: 1,
    paddingLeft: 20,
    paddingRight: 20
  }
}

module.exports = Preview;
