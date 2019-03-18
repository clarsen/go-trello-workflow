const express = require('express')
const path = require('path')
const dev = process.env.NODE_ENV !== 'production'
const next = require('next')
// const pathMatch = require('path-match')
const app = next({ dev })
const handle = app.getRequestHandler()
// const { parse } = require('url')

const server = express()
// const route = pathMatch()

server.use('/_next', express.static(path.join(__dirname, '.next')))

server.get('/', (req, res) => app.render(req, res, '/'))
// catch auth callback
server.get('/callback', (req, res) => app.render(req, res, '/callback'))
server.get('/login', (req, res) => app.render(req, res, '/login'))

// default
server.get('*', (req, res) => handle(req, res))

module.exports = server
