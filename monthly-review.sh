#!/bin/bash
out=$HOME/Dropbox/notes/personal/notational/monthly-retrospective.md

trello-dump-summary m \
  && generate-visualization tm \
  && generate-visualization m > $out && cat $out