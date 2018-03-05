#!/bin/bash

trello-dump-summary w && \
trello-dump-summary m && \
trello-dump-summary y \
  && generate-visualization yp
