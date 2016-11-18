/*
  window.mock.js

  This function is NOT tracked by flow, because it violates some flow
  rules in order to enable the correct use of the `window` object.

  This function mocks the Cookies object, which exists on the browser
  but not in the test environment.
*/

class CookiesObject {
  get(s: string, d?: any) {}
  set(s: string) {}
  remove(s: string) {}
}

// Add all necessary mock properties to the window object.
function transformWindow(w: any) {
  if(typeof w == "undefined") w = {};
  w.Cookies = new CookiesObject();
  w.document = {};
  w.addEventListener = (k: string, f: Function) => {};

  // Set the test URLs up.
  w.ORIGIN = "http://localhost:8080";
  w.HOST = "localhost";

  return w;
}

// This function returns the window object. If we're running on the browser,
// it'll return the browser's window object, otherwise it'll return a mock
// object that kind of works (for testing and prerendering on server side.)
function getWindowObject() {
  if(typeof window == "undefined") {
    return transformWindow({});
  } else {
    return window;
  }
}

module.exports = {
  transformWindow,
  getWindowObject
};
