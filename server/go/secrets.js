const dotenv = require("dotenv")

console.log("using .env." + process.env.STAGE)

// require and configure dotenv, will load vars in .env.$GQL_STAGE in PROCESS.ENV
dotenv.config({ path: ".env." + process.env.STAGE })

module.exports.getEnvVars = () => ({
  STAGE: process.env.STAGE,
  DOMAIN: process.env.DOMAIN,
  GITHUB_TOKEN: process.env.GITHUB_TOKEN,
  "AUTH0_CLIENT_ID": process.env.AUTH0_CLIENT_ID,
  appkey: process.env.appkey,
  authtoken: process.env.authtoken,
  user: process.env.user,
  USER_EMAIL: process.env.USER_EMAIL
})
