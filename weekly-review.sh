#!/bin/bash
out=$HOME/Dropbox/notes/personal/notational/retrospective.md
dump-summary && generate-visualization > $out && cat $out
