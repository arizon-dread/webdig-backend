# WebDig-backend
A rest api that takes a POST request to the /api/dig endpoint containting a json with a single field, "host" which in turn should contain either an IPv4 address or a dns-name. It will do a dnslookup and return a response with the result. 

Releases are provided as src code and an open dockerhub repo exists at [github](https://hub.docker.com/_/arizon/webdig-backend) [image at docker](docker.io/arizon/webdig-backend)

## Configuration
The DNS servers are configured using the confFile/config.yaml file.  
Example:

```yaml
---
general:
  cors:
    origins: ['http://localhost:4200']
    methods: ['GET', 'POST']
    headers: ['Origin', 'Content-Type']
dns:
  - name: "internet"
    servers: ['8.8.4.4', '8.8.8.8'] #ipv6 addresses must be enclosed in brackets []
    #This filters duplicate hits between different dns server group entries:
    filterDuplicates: false #if multiple dns server groups, at least one must be unfiltered (false). Duplicate hits within one dns group will still will always be filtered.
  - name: "OpenDNS"
    servers: ['208.67.222.222', '208.67.220.220']
    filterDuplicates: true #If true, hits in this group that exists in "internet" will be removed from this result section.
```
### filterDuplicates
The filterDuplicates setting is used if you specify several groups of dns servers and i.e. your internal dns servers mirror internet dns, but you want to display external addresses only in the 'internet' response. Then you would want to set `filterDuplicates: true` on your internal dns server group. The config here is just a public example. The point of using multiple dns server groups is to do dns lookups for different networks and avoid displaying mirrored hits in one or more groups.

Note that at least one dns server group must have `filterDuplicates: false`.  

Duplicate dns records or ip addresses within one servergroup will always be sorted and filtered, no matter this setting. If in this example above, `8.8.4.4` and `8.8.8.8` both reply with the same result for a lookup request, the response will only contain one entry of that hit for that dns server group. 

# Samples
The above config is used for these samples.
##  request 1
```json
{
    "host": "google.com"
}
```
## response 1
```json
{
   "error" : null,
   "results" : [
      {
         "dnsNames" : null,
         "error" : null,
         "ipAddresses" : [
            "142.250.74.110",
            "142.250.74.174"
         ],
         "name" : "internet"
      },
      {
         "dnsNames" : null,
         "error" : null,
         "ipAddresses" : [
            "142.250.74.142"
         ],
         "name" : "OpenDNS"
      }
   ]
}
```
## request 2
```json
{
    "host": "194.71.18.101"
}
```
## response 2
Note that `filterDuplicates: true` on OpenDNS filters all the hits here. If it would've been set to `false` the `dnsNames` arrays would've contained the same content. 
```json
{
   "error" : null,
   "results" : [
      {
         "dnsNames" : [
            "ica.se.",
            "mail2.ica.se.",
            "mail3.ica.se.",
            "mail4.ica.se.",
            "mobil.ica.se.",
            "static1.ica.se.",
            "static2.ica.se.",
            "static3.ica.se.",
            "www.ica.se."
         ],
         "error" : null,
         "ipAddresses" : null,
         "name" : "internet"
      },
      {
         "dnsNames" : [],
         "error" : null,
         "ipAddresses" : null,
         "name" : "OpenDNS"
      }
   ]
}
```
