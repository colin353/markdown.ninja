/*
  declarations.js
  @flow

  A place to put global flow declarations.
*/

declare class CookiesObject {
  get(s: string, d?: any) : void;
  set(s: string) : void;
  remove(s: string) : void
}

declare class Window {
  api: any,
  Cookies: any,
  document: any,
  location: any,
  _reactRoot: any,
  md5(target: string): string,
  ace: any,
  showdown: any,
  addEventListener: Function,
  open: (url: string) => void,
  innerWidth: number,
  innerHeight: number,
  fetch: Function,
  ORIGIN: string,
  HOST: string
};

type UploadEvent = {
  loaded: number,
  total: number
};

declare var window: Window;
