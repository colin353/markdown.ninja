/*
  api.test.js
  @flow

  Checks the API for normal behavior.
*/

// Flow setup for testing.
declare var test:any;
declare var expect:any;
declare var describe:any;
declare var jest:any;

var fetch = jest.fn(() => new Promise(resolve => resolve()));
var API = require('./api');

test("event subscriptions", () => {
  var api = new API.api();

  // Return a promise which will be fulfilled when the test is over.
  return new Promise((resolve, reject) => {
    api.addListener("clickBody", "testListener", (message) => {
      expect(message).toEqual("hello world");
      resolve();
    })
    api.emitMessage("clickBody", "hello world");
  });
});

test("event unsubscribe", () => {
  var api = new API.api();

  // Return a promise which will be fulfilled when the test is over.
  return new Promise((resolve, reject) => {
    api.addListener("clickBody", "testListener", (message) => {
      expect(message).toEqual("hello world");
      reject();
    })
    api.removeListeners("testListener");
    api.emitMessage("clickBody", "hello world");

    setTimeout(() => {
      resolve();
    }, 100)
  });
})
