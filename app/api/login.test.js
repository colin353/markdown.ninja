/*
  login.test.js
  @flow

  Test login, signup, and authentication functions.
*/

// Flow setup for testing.
declare var test:any;
declare var it:any;
declare var expect:any;
declare var describe:any;
declare var jest:any;

var fetch = jest.fn(() => new Promise(resolve => resolve()));
var API = require('./api');
var api = new API.api();

test("initial authentication should be not logged in", () => {
  return api.checkAuth().then((result) => {
    expect(result).toBe(false);
  });
})

test("log in with invalid credentials", () => {
  return api.login("test", "fake").then((result) => {
    expect(result).toBe(false);
  }).catch((err) => {
    expect(err).toBe("login failed");
  });
})

describe("sign up with a new account", () => {
  return api.signup({
    name: "Test User",
    domain: "testdomain",
    email: "test@test.com",
    password: "mytestpass"
  }).then((result) => {
    // Make sure no error occurred during signup.
    it("should not give an error when signing up", () => {
      expect(result.error).toBe(false);
    })
    // Now we should be automatically logged into the
    // account, so we'll run checkAuth
    return api.checkAuth();
  }).then((result) => {
    // Make sure we're logged in.
    if("should be logged in now, after signup", () => {
      expect(result).toBe(true);

      // Make sure that the user data is correct.
      expect(api.user.name).toBe("Test User");
      expect(api.user.domain).toBe("testdomain");
      // Make sure we don't accidlentally get some weird fields
      // that shouldn't be there.
      expect(api.user.password).toBeUndefined();
      expect(api.user.password_hash).toBeUndefined();
      expect(api.user.password_salt).toBeUndefined();
    });

  });
})
