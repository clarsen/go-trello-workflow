#!/bin/bash

trello-dump-summary w \
  && generate-visualization tw
(find templates -type f; \
 find ~/lsrc/data-and-reviews/reviews -type f; \
 find ~/lsrc/data-and-reviews/task-summary -type f) \
| entr generate-visualization w
