# dynamic-headers
dynamic caddy server headers middleware. 


# build Caddy with dynamic headers module
 
xcaddy build --output ./caddy.exe --with github.com/c1pca/dynamic-headers

```
{
  "logging": {
    "logs": {
      "default": {
        "level": "DEBUG"
      }
    }
  },
  "apps": {
    "http": {
      "servers": {
        "myserver": {
          "listen": [":448",":87"],
          "routes": [
            {
              "match": [{"path": ["/*"]}],
              "handle": [
                {
                  "handler": "dynamic_headers",
                  "to_header": "test_dynamic_header_1",
                  "from_header": "Connection", // optional parameter
                  "take_host": true
                },
                {
                "handler": "reverse_proxy",
                "headers": {
                  "request": {
                    "set": {
                      "X-Static-Header": [
                        "test-static-Header"
                      ]
                    }
                  }
                },
                "upstreams": [{
                  "dial": "localhost:445"
                }]
                }
              ]
            }
          ],
          "logs": {
            "should_log_credentials": true
          }
        }
      }
    }
  }
}
```