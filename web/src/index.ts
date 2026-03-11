import { Container } from './interfaces'

const eSource = new EventSource('http://localhost:6060/containers')

eSource.addEventListener('message', (e) => {
  const containers: Container[] = JSON.parse(e.data)

  console.log(containers)
})
