/*
  debounce.js
  @flow

  Takes a function and returns a debounced version of the function,
  which is called only once every $delay milliseconds even if the
  original function is called much more often.
*/

module.exports = function(fn: Function, delay: number) : Function {
  var timer: number = 0;
  return function () {
    var context = this, args = arguments;
    clearTimeout(timer);
    timer = setTimeout(function () {
      fn.apply(context, args);
    }, delay);
  };
}
