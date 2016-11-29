/*
  indexcss.js

  This function builds a definition file for the custom CSS styles.
*/

var fs = require('fs');
var glob = require('glob')

// Grab a list of the css files under the
// webstyles directory.
glob("../web/css/webstyles/*.css", {}, (err, files) => {
  var available_styles = [];

  for(const f of files) {
    var filename = f.replace(/^.*[\\\/]/, '').slice(0,-4);
    if(filename == "required") continue;

    // Add that filename to the list.
    available_styles.push(filename);
  }

  // Also have to write it in the npm-accessible directory so it can be read
  // by node.
  fs.writeFile("./config/styles.json", JSON.stringify({styles: available_styles}), (err) => {
    if(err) {
      console.log("Err: couldn't write style list file!");
      throw err;
    } else {
      console.log("Wrote styles.json.")
    }
  });
});
