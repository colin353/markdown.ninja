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

type User = {
  name: string,
  domain: string,
  external_domain: string,
  email: string
}
