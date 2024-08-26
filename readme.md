# 2wayssl
A Two-Way ssl generator template

- [x] Auto generated two-way ssl certificate
- [x] Auto generated nginx template
- [x] An simple 2way server/client example in golang

## install
    make install
## USAGE
Usage:
    
    $ 2wayssl -h
    2wayssl [-p PORT] [--silent] -d your-domain.com

Example:

    $ 2wayssl - 2wayssl -d 2wayssl.local

    echo "127.0.0.1 2wayssl.local" | sudo tee -a /etc/hosts
    cd ~/.2wayssl && curl --cacert ca.crt --cert  client.crt --key client.key --tlsv1.2  https://2wayssl.local:444

                    +---------------------------+
                    |curl -k https://local1.com |
                    +------+--------------------+
                              |
                              v 
                      +-------+------+
                      | nginx gateway| default port: 443
                      | (port:444)   |  
                      ++-----+-------+  
                         |         | 
                         v         v
               +-------+---+        +-----------+  
               | upstream1 |        | upstream2 |  
               |(port:4500)|        |(port:4501)|  
               +-----------+        +-----------+  
                   
