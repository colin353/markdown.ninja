/*
  build.js

  This program runs some general build commands.
*/

var fs = require('fs');
var glob = require('glob')
var showdown = require('showdown');

var converter = new showdown.Converter();

// Take the markdown files and turn them into html files in the
// "./web/default" folder.
glob("../web/default/*.md", {}, (err, files) => {
  for(const f of files) {
    console.log("Rendering: ", f, " --> HTML");
    
    fs.readFile(f, 'utf-8', (err, data) => {
      if(err) {
        console.log("Err: couldn't read file! ", f)
        throw err;
      }
      fs.writeFile(f+".html", converter.makeHtml(data), (err) => {
        if(err) {
          console.log("Err: couldn't write file!", f+".html");
          throw err;
        }
      });
    })
  }
});
