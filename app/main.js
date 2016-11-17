/*
  main.js
  @flow
*/

var ReactDOM = require('react-dom');
var Routes   = require('./routes');

window._reactRoot = document.getElementById('root');

ReactDOM.render(
  Routes,
  window._reactRoot
);
