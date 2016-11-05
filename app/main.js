/*
  main.js
  @flow
*/

var ReactDOM = require('react-dom');
var Routes   = require('./routes');

window._reactRoot = document.getElementById('root');

var API = require('./api/api');

window.api = new API.api();
window.api.checkAuth();

ReactDOM.render(
  Routes,
  window._reactRoot
);
