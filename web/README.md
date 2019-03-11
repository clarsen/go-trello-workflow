# prepare AWS certificate
- go to AWS certificate manager
- Request Certificate -> public certificate
- use domain in .env.production
- DNS verification
- Create record in Route 53

# prepare auth0
- create single page app
- update `lib/auth0.js` AUTH0_CLIENT_ID, AUTH0_DOMAIN
- update `serverless/.env.production` AUTH0_CLIENT_ID
- under advanced settings, copy signing certificate to `serverless/python-auth0/public_key`
- add to allowed web origins https://enchilada-serverless-next-auth0.app.caselarsen.com
- add to allowed callback urls https://enchilada-serverless-next-auth0.app.caselarsen.com/callback
- add to allowed logout URLs https://enchilada-serverless-next-auth0.app.caselarsen.com

# running locally
```
npm run dev
```

# domain (first time before deploy)
```
NODE_ENV=production sls create_domain
```

# deploying
```
NODE_ENV=production npm run deploy
```

# test
- go to https://enchilada-serverless-next-auth0.app.caselarsen.com
  - click "try it" without logging in, see that we get an error
  - log in and then click "try it", see that we get an "ok" response
