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
var AuthenticationWrapper = require('./components/authenticationwrapper');

var routes = (
  <Router history={ReactRouter.browserHistory}>

    <Route path='/edit/site' component={FullScreenWrapper}>
      <Route component={AuthenticationWrapper}>
        <IndexRoute component={require('./pages/edit')} />
      </Route>
    </Route>
    <Route path='/' component={Wrapper}>
      <IndexRoute component={require('./pages/index')} />
      <Route path='/edit/login' component={require('./pages/login')} />
      <Route path='/edit/signup' component={require('./pages/signup')} />
      <Route path="*" component={require('./pages/404')}/>
    </Route>
  </Router>
);

module.exports = routes;
