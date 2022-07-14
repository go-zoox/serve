#!/bin/sh

serve \
  -p ${PORT:-8080} \
  -d ${DIR:-/var/www/html} \
  --prefix ${PREFIX:-/}
