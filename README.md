# go-trello-workflow

## Installation
    $ go get github.com/clarsen/go-trello-workflow
    $ trello-workflow-cli



    generate Trello app key by https://trello.com/1/appKey/generate
    get auth token from https://trello.com/1/connect?key=<YOUR TRELLO APP KEY>&name=trellow-workflow&response_type=token

Trello boards should be set up as per expectations of code.

## Deploy to heroku
$ heroku config:set appkey=<YOUR TRELLO APP KEY>
$ heroku config:set authtoken=<YOUR TRELLO ACCOUNT AUTH TOKEN>
$ heroku config:set user=<your trello username>
$ git push heroku master

Turn on `longrun` dyno which runs go-trello-workflow

## Use in CLI
Create .env with:
    appkey=<YOUR TRELLO APP KEY>
    authtoken=<YOUR TRELLO ACCOUNT AUTH TOKEN>
    user=<your trello username>

Then, these commands:

    $ trello-workflow-cli -h

    NAME:
       trello-workflow - Update Trello board

    USAGE:
       trello-workflow-cli [global options] command [command options] [arguments...]

    VERSION:
       0.0.0

    COMMANDS:
         today, t  Update the trello board on daily basis
         help, h   Shows a list of commands or help for one command

    GLOBAL OPTIONS:
       --help, -h     show help
       --version, -v  print the version

### In the morning or end of day
    $ trello-workflow-cli today

## Notable mentions
- Butlerbot does this using trello cards as a conversational API https://trello.com/b/2dLsEE9t/butler-for-trello
  and has a command builder https://butlerfortrello.com/builder.html   I will use this for the orange label cherry picking interaction.

## License

MIT
