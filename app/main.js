/*
  main.js
  @flow
*/

var ReactDOM = require('react-dom');
var Routes   = require('./routes');

window._reactRoot = document.getElementById('root');
//window.api = new APIType();

ReactDOM.render(
  Routes,
  window._reactRoot
);
