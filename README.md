# go-trello-workflow

## Installation
```
    git clone https://github.com/clarsen/go-trello-workflow.git
```


Generate Trello app key by https://trello.com/1/appKey/generate

Get auth token from https://trello.com/1/connect?key=<YOUR TRELLO APP KEY>&name=trellow-workflow&response_type=token

Trello boards should be set up as per expectations of the code.

## Installation (for dev)
```
    git clone https://github.com/clarsen/go-trello-workflow.git
    cd web
    npm install
```

## Deploy to AWS serverless

- automatic on git push via CircleCI
- ensure envrionment variables with API keys, email address are set up.

## Notable mentions
- Butlerbot does this using trello cards as a conversational API https://trello.com/b/2dLsEE9t/butler-for-trello
  and has a command builder https://butlerfortrello.com/builder.html   I will use this for the orange label cherry picking interaction.


## License

MIT
