console.log('Hello, world')

const eSource = new EventSource('http://localhost:6060/containers')

eSource.addEventListener('message', (e) => {
  console.log(e.data)
})
