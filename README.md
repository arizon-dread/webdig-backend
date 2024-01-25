# WebDig-backend
A rest api that takes a POST request to the /api/dig endpoint containting a json with a single field, "host" which in turn should contain either an IPv4 address or a dns-name. It will do a dnslookup and return a response with the result. 

Releases are provided as src code and an open dockerhub repo exists at (https://hub.docker.com/_/arizon/webdig-backend)[docker.io/arizon/webdig-backend]
