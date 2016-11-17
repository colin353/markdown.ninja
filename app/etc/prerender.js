var React = require('react');
var ReactDOMServer = require('react-dom/server');

var { Wrapper, FullScreenWrapper} = require('../components/wrappers');

var routes = require('./routes.json');

console.log(routes);

var API = require('../api/api');
var api = new API.api();

for(var route in routes) {
  var Component = require('../pages/' + routes[route].component);
  var output = ReactDOMServer.renderToStaticMarkup(
    <Wrapper api={api}>
      <Component api={api} />
    </Wrapper>
  );
  console.log(output);
}


console.log("done");
