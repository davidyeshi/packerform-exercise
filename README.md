# packerform-exercise

## Cloning the package
- don't forget the frontend code, which would require a separate clone

## System Requirements

- mongo server running
- GOPATH set on the right directory

### Importing data to mongo
- cd to ./mongo-script
- command: go run loaddata.go

**Note: mongo-server port is set to default:27017 .To change this go to loaddata.go/getMongoClient()

### Running backend go server
- cd to ./backend-golang
- command: go run main.go mongohelper.go structs.go
- the server will run on http://localhost:8000

**Note: To change mongo-server port in backend, go to mongoHelper.go/getMongoClient()

### Running frontend
So once backend is running
- cd to ./frontend-react
- do a npm install
- npm start
- the server will run on http://localhost:3000

#### Cheers!
