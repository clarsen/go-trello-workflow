const dotenv = require("dotenv")

console.log("using .env." + process.env.NODE_ENV)

// require and configure dotenv, will load vars in .env.$GQL_STAGE in PROCESS.ENV
dotenv.config({ path: ".env." + process.env.NODE_ENV })

module.exports.getEnvVars = () => ({
  NODE_ENV: process.env.NODE_ENV,
  DOMAIN: process.env.DOMAIN
})
