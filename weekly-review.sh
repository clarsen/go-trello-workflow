#!/bin/bash
out=$HOME/Dropbox/notes/personal/notational/retrospective.md
trello-dump-summary w \
  && generate-visualization tw \
  && generate-visualization w > $out && cat $out
