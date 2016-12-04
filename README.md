## markdown.ninja

![build-status](https://circleci.com/gh/colin353/markdown.ninja.png?style=shield)

This is a little project to make a service which helps people maintain
a personal portfolio website, which shows off their personal projects.

## Technology

The web server is written in Go, and the frontend uses React. Continuous integration/deployment
is provided by [by CircleCI](https://circleci.com/gh/colin353/markdown.ninja). You can see build status for previous builds on that page, and the configuration file for running the builds/deployment in `/circle.yml`.

I'm running the site ([http://markdown.ninja](http://markdown.ninja)) on Google's Container Engine in a three node kubernetes cluster. Since the system has to host and serve files, those are served by [GlusterFS](https://www.gluster.org/). The database of user data is stored in [Redis](https://redis.io/).

Docker images for the project (only about 10 MB in size!) are hosted here: https://hub.docker.com/r/colinmerkel/portfolio (that was the old name of this project). The docker images are so small because I'm using `FROM scratch` in my `Dockerfile`, and a statically linked go binary as the server.

## Contributing guidelines

Feel free to send a PR if you want a feature added/fixed! One thing I would like in
particular is new styles. If you want to contribute new styles, you can do so by:

1. Creating a file under /web/css/webstyles of the form `<STYLE NAME>.css`.
2. All of the rules in the stylesheet should be prefixed with the `.md_container` class. Take a look at the existing stylesheets to get an idea of what I mean by this.
3. Make sure your stylesheets look sensible on both mobile and desktop. 

Also, if you want to contribute an example of a site made using markdown.ninja for the front page,
you can add a screenshot of your site under `/web/img/`.

## Support this site

Using markdown.ninja is free for everyone. But it costs money to host the service. If you want to
contribute a few dollars toward the hosting costs of the site, [consider doing so using PayPal](https://www.paypal.me/markdownninja). All of the donations go toward hosting costs. Thanks!
