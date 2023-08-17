# hello-websockets

Dummy http server with websocket support based on [chi](https://github.com/go-chi/chi) and [nhooyr/websocket](https://github.com/nhooyr/websocket).

## Usage
```sh
docker image build -t hello-ws .
docker container run --rm -d -p 8080:8080 --name hello-ws hello-ws
curl localhost:8080
# {"message":"Hello!"}
```

### Test html page
Appends all messages received to table.

```
<!DOCTYPE html>
<html lang="en">
  <body>
    <h2>Hello World</h2>

    <table id="my-table">
      <thead>
        <tr>
          <th>Messages Received</th>
        </tr>
      </thead>
      <tbody>
      </tbody>
    </table>

    <script>
        let socket = new WebSocket('ws://127.0.0.1:8080/ws');
        let tableBody = document.getElementById('my-table').getElementsByTagName('tbody')[0];
        console.log('Connecting...');

        socket.onopen = () => {
            console.log('Successfully Connected');
            socket.send('{"message": "Just Connected"}');
        };

        socket.onclose = event => {
            console.log('Socket Closed Connection: ', event);
        };

        socket.onerror = error => {
            console.log('Socket Error: ', error);
        };

        socket.addEventListener("message", ({ data }) => {
          let row = tableBody.insertRow();
          let cell = row.insertCell();
          cell.appendChild(document.createTextNode(data))
        });

        setInterval(() => {
            let message = 'Current timestamp: ' + Date.now() ;
            socket.send('{"message": "' + message + '"}');
        }, 2000)
    </script>
  </body>
</html>
```