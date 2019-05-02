import Router from 'next/router'
// from with-apollo-auth example

export default (context, target) => {
  if (context.res) {
    // server
    // 303: "See other"
    console.log('303 -> ', target)
    context.res.writeHead(303, { Location: target })
    context.res.end()
  } else {
    // In the browser, we just pretend like this never even happened ;)
    console.log('Router replace -> ', target)
    Router.replace(target)
  }
}
