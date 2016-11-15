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

var { Wrapper, FullScreenWrapper } = require('./components/wrappers');

var routes = (
  <Router history={ReactRouter.browserHistory}>
    <Route path='/edit/site' component={FullScreenWrapper}>
      <IndexRoute component={require('./pages/edit')} />
    </Route>
    <Route path='/' component={Wrapper}>
      <Route path='/edit/login' component={require('./pages/login')} />
      <Route path='/edit/signup' component={require('./pages/signup')} />
      <Route path="*" component={require('./pages/404')}/>
    </Route>
  </Router>
);

module.exports = routes;
