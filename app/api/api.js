/*
  api.js
  @flow

  This file defines the interface for communicating with the webserver.
*/

type User = {
  name: string,
  phone_number: string,
  domain: string,
  email: string
}

class API {
  BASE_URL: string;
  BASE_DOMAIN: string;
  callbacks: { [key: string]: {[key: string]: Array<Function>}};
  authenticated: bool;
  user: User;

  constructor() {
    // Event callback dictionary.
    this.callbacks = {
      authenticationStateChanged: {},
      escapeKeyPressed: {},
      clickBody: {},
      "ctrl+s": {}
    };

    // The base URL is stored here so that UI components can
    // get access to it via window.api.BASE_URL.
    if(window) {
      this.BASE_URL = window.location.origin;
      this.BASE_DOMAIN = window.location.host;
    }

    // Check the cookies for a quick sense of whether
    // we are currently logged in.
    this.authenticated = false;
    if(window.Cookies && window.Cookies.get("authenticated"))
      this.authenticated = true;

    // Detect the escape key, which is used for quitting dialogs
    // and things like that.
    window.document.onkeyup = (e) => {
       if (e.keyCode == 27) {
         this.emitMessage('escapeKeyPressed');
      }
    };

    // Detect when the body of the website is clicked, which is used
    // for cancelling menus and things of that nature.
    if(window._reactRoot)
      window._reactRoot.addEventListener('click', (event) => {
        this.emitMessage('clickBody', event);
      });

    // Detect CTRL+S, block propagation, and send a message out about it.
    window.addEventListener("keydown", (e) => {
      if((e.ctrlKey || e.metaKey) && (e.which == '115' || e.which == '83')) {
        event.preventDefault();
        this.emitMessage("ctrl+s", e);
      }
      return false;
    });
  }

  // The basic function to talk to the API. It encodes the parameters as JSON
  // and sends them as a POST request.
  request(uri: string, params: any) {
    var request_url = this.BASE_URL + uri;
    return fetch(request_url, {
        'headers'     : {'Content-Type': 'application/json'},
        'credentials' : 'include',
        'method'      : 'POST',
        'body'        : JSON.stringify(params)
    }).then((r) => {
        if(r.status != 200) return r.json().then((z) => { throw z;});
        else return r.json();
    });
  }

  setAuthenticationState(state: boolean) {
    this.authenticated = state;
    if(state)
      window.Cookies.set("authenticated", true);
    else window.Cookies.remove("authenticated");
    this.emitMessage('authenticationStateChanged', this.authenticated);
  }

  login(domain: string, password: string) {
    return this.request("/api/auth/login", {domain, password}).then((response) => {
      if(response.result == "ok")
        return this.checkAuth(true);
      else return response;
    });
  }

  logout() {
    return this.request("/api/auth/logout").catch(() => {}).then(() => {
      return this.checkAuth(false);
    });
  }

  checkAuth() {
    return this.request("/api/auth/check").then((response) => {
      this.user = response.user;
      this.setAuthenticationState(true);
      return true;
    }).catch(() => {
      this.setAuthenticationState(false);
      return false;
    });
  }

  emitMessage(message: string, data: any) {
    // Don't do anything if nobody is listening to the message.
    if(this.callbacks[message] === undefined) return;
    for (var key in this.callbacks[message]) {
      for (var i=0; i<this.callbacks[message][key].length; i++) {
        this.callbacks[message][key][i](data);
      }
    }
  }

  addListener(message: string, key: string, callback: Function) {
      if(this.callbacks[message] == undefined) {
        console.warn("Registered for callback", message, ", which doesn't exist.");
        this.callbacks[message] = {};
      }
      // Check if the key has never registered callbacks.
      if (!(key in this.callbacks[message])) this.callbacks[message][key] = [];
      this.callbacks[message][key].push(callback);
  }

  // When cleaning up, call this method with the listener key, so we can
  // properly delete the listeners. It should be called in willUnmount, I guess.
  removeListeners(key: string) {
    for (var message in this.callbacks) {
      if (key in this.callbacks[message])
        delete this.callbacks[message][key];
    }
  }

  signup(params: {name: string, domain: string, password: string}) {
    return this.request("/api/auth/signup", params);
  }

  pages() : Promise<Page[]> {
    return this.request("/api/edit/pages");
  }

  createPage(params: Page) {
    return this.request("/api/edit/create_page", params);
  }

  editPage(params: Page) {
    return this.request("/api/edit/edit_page", params);
  }

}

type APIInstance = API;
export type { APIInstance };

module.exports = {
  api: API
};
