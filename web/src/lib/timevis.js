import { ENDPOINTS } from './api'
import auth from './auth0'
import fetch from 'isomorphic-unfetch'

const fetchTimeReport = () => {
  console.log('fetching time report')
  return fetch(ENDPOINTS['python']['timevis_api'], {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${auth.instance().getIdToken()}`,
    },
  })
    .then(response => {
      console.log('got', response)
      return response.json()
    })
//     .then((data) => {
//       console.log('Token:', data);
//       this.setState({
//         content: data.message,
//         loading: false
//       })
//       // document.getElementById('message').textContent = '';
//       // document.getElementById('message').innerHTML = data.message;
//     })
}

export default fetchTimeReport
