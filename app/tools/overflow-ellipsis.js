/*
  overflow-ellipsis.js
  @flow

  Returns a string limited to a fixed number of characters.
*/

function overflowEllipsis(text: string, limit: number) : string {
  if(text.length <= limit) return text;
  else return text.slice(0, limit-3) + "...";
}

module.exports = overflowEllipsis;
