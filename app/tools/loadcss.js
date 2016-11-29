/*
  loadcss.js
  @flow

  The purpose of this file is to allow dynamically loading and
  unloading a CSS file.
*/

module.exports = {
  load: (path: string) => {
    var headNode = document.getElementsByTagName("head")[0];
    var cssNode = document.createElement('link');
    cssNode.type = 'text/css';
    cssNode.rel = 'stylesheet';
    cssNode.href = path;
    cssNode.media = 'screen';
    cssNode.setAttribute("data-dynamic", "true");
    headNode.appendChild(cssNode);
  },

  unload: () => {
    var targetElement = "link";
    var targetAttr = "href";

    var links = document.getElementsByTagName(targetElement);
    for (var i=links.length; i>=0; i--)  {
      if (links[i] && links[i].getAttribute('data-dynamic')=="true") links[i].parentNode.removeChild(links[i]);
    }
  }
}
