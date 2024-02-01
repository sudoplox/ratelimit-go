# ratelimit-go

## To run the server

```
go run .
```

## To test token-bucket

```
$ for i in {1..6}; do curl http://localhost:8080/ping; done

 {"status":"Successful","body":"Hi youve reached the ApAPI"}
 {"status":"Successful","body":"Hi youve reached the ApAPI"}
 {"status":"Successful","body":"Hi youve reached the ApAPI"}
 {"status":"Successful","body":"Hi youve reached the ApAPI"}
 {"status":"Request Failed","body":"The API is at capacity, please try again later."}
 {"status":"Request Failed","body":"The API is at capacity, please try again later."}
```