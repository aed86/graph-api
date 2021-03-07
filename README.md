##### To start project please do:

Run neo4j: `docker-compose up neo4j` 

Apply fixtures: `go run bin/fixtures_apply.go`

Run app: `go run main.go`



####In the project I've used graph database Neo4j and golang backend.
I have realised base REST API with main handlers which allow get, add, update and delete graph nodes and add, delete relations between these nodes.

### How to use:

###### To find shortest path between nodes run:

`curl --location --request POST 'localhost:3000s' \
--header 'Content-Type: application/json' \
--data-raw '{
"sourceName": "A",
"targetName": "F"
}'`

###### Get all nodes

`curl --location --request GET 'localhost:3000/' \
--header 'Content-Type: application/json' \
--data-raw '{
"limit": 100
}'`

###### Get node by id

`curl --location --request GET 'localhost:3000/node/40' \
--header 'Content-Type: application/json'`

###### Add node

`curl --location --request POST 'localhost:3000/node' \
--header 'Content-Type: application/json' \
--data-raw '{
"name": "John7",
"born": 123
}'`

###### Add relation

`curl --location --request POST 'localhost:3000/relation' \
--header 'Content-Type: application/json' \
--data-raw '{
"source": 13,
"target": 12,
"cost": 10
}'`

###### Delete relation

`curl --location --request DELETE 'localhost:3000/relation' \
--header 'Content-Type: application/json' \
--data-raw '{
"target": 3,
"source": 0
}'`

###### Get all neighbours for node

`curl --location --request GET 'localhost:3000/node/52/neighbours' \
--header 'Content-Type: application/json' \
--data-raw '{
"limit": 10
}'`


###### Remove node by id
`curl --location --request DELETE 'localhost:3000/node/8' \
--header 'Content-Type: application/json' \
--data-raw ''`