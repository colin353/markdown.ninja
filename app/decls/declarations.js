/*
  declarations.js
  @flow

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
  md5(target: string): string,
  ace: any,
  showdown: any,
  addEventListener: Function,
  open: (url: string) => void,
  innerWidth: number,
  innerHeight: number
};

type UploadEvent = {
  loaded: number,
  total: number
};

declare var window: Window;
