# prepare AWS certificate
- go to AWS certificate manager
- Request Certificate -> public certificate
- use domain in .env.production
- DNS verification
- Create record in Route 53
- wait up to 40 minutes or so, however the certificate usually becomes available sooner than that.

# domain (first time before deploy)
```
npm install
STAGE=production sls create_domain
```

# deploy manually

```
    export AWS_PROFILE=serverlesstest
    export STAGE=production
    npm install
    make deploy
```


# test

```
  echo '{"action": "morning-reminder"}' | sls invoke -f scheduled -l
  echo '{"action": "test"}' | sls invoke -f scheduled -l
```
