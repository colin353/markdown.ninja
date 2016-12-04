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

var API = require('./api/api');
var api = new API.api();

// Add the api to the window object, so debugging is a bit
// easier.
window.api = api;
api.checkAuth();

// This is necessary to pass the api prop down to children.
const useExtraProps = {
  renderRouteComponent: child => React.cloneElement(child, {api})
};

var routes = (
  <Router render={ReactRouter.applyRouterMiddleware(useExtraProps)} history={ReactRouter.browserHistory}>

    <Route path='/edit/site' component={FullScreenWrapper}>
      <Route component={AuthenticationWrapper}>
        <IndexRoute component={require('./pages/edit')} />
      </Route>
    </Route>
    <Route path='/' component={Wrapper}>
      <IndexRoute component={require('./pages/index')} />
      <Route path='/edit/login' component={require('./pages/login')} />
      <Route path='/edit/signup' component={require('./pages/signup')} />
      <Route path='/edit/account' component={require('./pages/account')} />
      <Route path="*" component={require('./pages/404')}/>
    </Route>
  </Router>
);

module.exports = routes;
