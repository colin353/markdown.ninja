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
      <div className="md_container" style={styles.container} dangerouslySetInnerHTML={{ __html: this.props.html}}></div>
    );
  }
}

const styles = {
  container: {
    flex: 1,
    zIndex: 4,
    boxShadow: 'inset 0px 0px 5px #000',
    paddingLeft: 20,
    paddingRight: 20,
    overflow: 'auto'
  }
}

module.exports = Preview;
