/*
  main.js
  @flow
*/

var ReactDOM = require('react-dom');

window._reactRoot = document.getElementById('root');

var Routes   = require('./routes');

ReactDOM.render(
  Routes,
  window._reactRoot
);
