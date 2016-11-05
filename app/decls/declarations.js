/*
  declarations.js

  A place to put global flow declarations.
*/

declare class Window {
  api: any,
  Cookies: {
    get: Function,
    set: Function,
    remove: Function
  },
  document: any,
  location: any,
  _reactRoot: any,
  md5(target: string): string
};

declare var window: Window;
