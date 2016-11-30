FROM scratch
MAINTAINER Colin Merkel <colinmerkel@podkarma.com>

# Add the application binary.
ADD server /

# Add the supplementary files.
ADD web /web
ADD config /config

# Create the data directory, just in case
# we want to use local storage.
ADD data /data

CMD ["/server"]

EXPOSE 8080
