/*
  types.js
  @flow

  Definitions for global types which are returned by the API.
*/

declare type Page = {
  name: string,
  markdown: string,
  html?: string
}

declare type File = {
  name: string
}
