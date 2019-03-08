localbuild: clean
	(cd cmds/trello-workflow-cli; go build -o ../../bin/trello-workflow-cli)
	(cd cmds/generate-visualization; go build -o ../../bin/generate-visualzation)
	(cd cmds/trello-dump-summary; go build -o ../../bin/trello-dump-summary)


.PHONY: clean
clean:
	rm -rf ./bin
