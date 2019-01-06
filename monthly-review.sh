#!/bin/bash

trello-dump-summary m \
  && generate-visualization tm
generate-visualization mri

(find templates -type f; \
 find ~/lsrc/data-and-reviews/reviews -type f; \
 find ~/lsrc/data-and-reviews/task-summary -type f) \
| entr generate-visualization m
