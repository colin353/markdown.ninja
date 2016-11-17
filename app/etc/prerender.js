/*
  prerender.js

  This script uses server-side rendering to generate HTML files
  for common routes. That way, when the browser goes to those places
  it can get prerendered HTML, making the user experience seem a bit
  faster.
*/

var React = require('react');
var ReactDOMServer = require('react-dom/server');

var fs = require('fs');

var { Wrapper, FullScreenWrapper} = require('../components/wrappers');

type Route = {
  file: string,
  component: string
}

// The routes are listed in this routes.json file, which contains a list
// of components we should try to prerender. This same json file is read
// by the server when it starts up, so it knows which routes to provide
// prerendered HTML for.
var routes : { (key:string) : Route } = require('../../config/routes.json');

// We'll need to create a copy of the API (even though it won't work
// properly in this server-side environment, it is needed for the
// functionality of some of the components we'll be rendering).
var API = require('../api/api');
var api = new API.api();

// The file ./web/index.html is the base HTML file, and we'll be inserting
// the prerendered content into the <div id="root" /> (react root component).
// The component is labelled with a comment like this:
//    <!-- react root -->
// so we'll search for that comment and then dump the rendered component in
// there.
new Promise((resolve, reject) => {
  fs.readFile("../web/index.html", 'utf-8', (err, data) => {
    if(err) reject("failed to load index.html");
    else resolve(data);
  });
}).then((data) => {
  // We're looking for a special comment like: <!-- react root -->.
  var token = "<!-- react root -->";
  var index = data.indexOf(token);
  if(index == -1) throw "unable to find identifier comment in ./web/index.html";

  var _beforeSplit = data.slice(0, index);
  var _afterSplit = data.slice(index + token.length);

  // We're going to create a function that takes the generated HTML for
  // a component, and returns the HTML that the prerendered file should contain
  // by wrapping it with the contents of index.html.
  function wrapHTML(html) {
    return _beforeSplit + html + _afterSplit;
  }

  return wrapHTML;
}).then((wrapHTML) => {
  // Loop through the components listed in the routes.json file,
  // and try to prerender each one.
  for(var route in routes) {
    var Component = require('../pages/' + routes[route].component);

    console.log("Rendering: ", route, " --> ", routes[route].file);
    var componentHTML = ReactDOMServer.renderToStaticMarkup(
      <Wrapper api={api}>
        <Component api={api} />
      </Wrapper>
    );

    // Wrap the component HTML with the contents of index.html
    var output = wrapHTML(componentHTML)

    // Write the contents of the file to disk.
    fs.writeFile("../web/prerendered/" + routes[route].file, output, (err) => {
      if(err) {
        console.log("Err: couldn't write file!", routes[route].file+".html");
        throw err;
      }
    });
  }
}).catch((err) => {
  console.log(err);
  console.log("Fatal error, terminating.");
  process.exit(1);
})

console.log("done");
