/*
  cookies.mock.js
  @flow

  This function mocks the Cookies object, which exists on the browser
  but not in the test environment.
*/

class CookiesObject {
  get(s: string, d?: any) {}
  set(s: string) {}
  remove(s: string) {}
}

// Add all necessary mock properties to the window object.
function transformWindow(w: Window) {
  if(typeof w == "undefined") w = {};
  w.Cookies = new CookiesObject();
  w.document = {};
  w.addEventListener = (k: string, f: Function) => {};

  // Set the test URLs up.
  w.ORIGIN = "http://localhost:8080";
  w.HOST = "localhost";
}

module.exports = transformWindow;
