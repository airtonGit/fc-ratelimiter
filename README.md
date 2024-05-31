
# Rate Limiter v0

## Configuring Rate Limiter with Environment Variables and .env file
This Rate Limiter service uses environment variables to configure its behavior. Here's how you can configure it:
### IP Limit:
   Variable: IP_LIMIT_SEC
   Description: Sets the maximum number of requests allowed per IP address per second.
   Default Value: 1
   Example: IP_LIMIT_SEC=10 will allow 10 requests per IP address per second.
### Token Limits:
   
   Variable: TOKENS_LIMIT_SEC

   Description: 
   Sets the maximum number of requests allowed per token per second. This is a comma-separated list of token-limit pairs.
   
   Default Value: sampleToken1=1,sampleToken2=3
   
   Example: TOKENS_LIMIT_SEC=myToken=5,anotherToken=10 
   
   will allow 5 requests per second for myToken and 10 requests per second for anotherToken.
### Redis Host:
   Variable: REDIS_HOST

   Description: Specifies the hostname and port of the Redis server.
   
   Default Value: localhost:6379
   
   Example: REDIS_HOST=my-redis-server:6379
### Redis Database:
   Variable: REDIS_DB
   
   Description: Specifies the Redis database to use.
   
   Default Value: 0
   
   Example: REDIS_DB=1

## Example .env file:
```
IP_LIMIT_SEC=10
TOKENS_LIMIT_SEC=myToken=5,anotherToken=10
REDIS_HOST=my-redis-server:6379
REDIS_DB=1
```