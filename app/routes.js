/*
  routes.js
  @flow

  This file defines the URL routing within the app.
*/

var React = require('react');

var ReactRouter = require('react-router');

var Router = ReactRouter.Router;
var Route = ReactRouter.Route;
var IndexRoute = ReactRouter.IndexRoute;

var routes = (
  <Router history={ReactRouter.browserHistory}>
    <Route path='/login' component={require('./pages/login')} />
    <Route path="*" component={require('./pages/404')}/>
  </Router>
);

module.exports = routes;
