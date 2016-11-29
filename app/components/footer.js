/*
  footer.js
  @flow

  This renders the bottom footer of the site.
*/

var React = require('react');

class Footer extends React.Component {
  render() {
    return (
      <div style={styles.container}>
        <p style={styles.footer}>Written by <a href="http://colin.markdown.ninja">Colin Merkel</a>, 2016. Source code at <a href="https://github.com/colin353/portfolio">https://github.com/colin353/portfolio</a>.</p>
      </div>
    );
  }
}

const styles = {
  container: {
    marginLeft: 20,
    marginRight: 20,
    display: 'flex',
    justifyContent: 'center'
  },
  footer: {
    fontSize: 12
  }
}

module.exports = Footer;
