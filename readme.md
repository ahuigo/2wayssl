# 2wayssl
Automatically generate Two-Way ssl certificates, configuration, and test server.

- [x] Auto generated two-way ssl certificate
- [x] Auto generated nginx template
- [x] Simple 2way-ssl server/client example in golang


## install
    git clone https://github.com/ahuigo/2wayssl
    cd 2wayssl && make install

## USAGE
Usage:
    
    $ 2wayssl -h
    2wayssl [-p PORT] [--silent] -d your-domain.com

Example:

    # 1. generate certificate and start a test server
    $ 2wayssl -p 444 -d 2wayssl.local

    # 2. test certificate via test server
    $ echo "127.0.0.1 2wayssl.local" | sudo tee -a /etc/hosts
    $ cd ~/.2wayssl && curl --cacert ca.crt --cert  client.crt --key client.key --tlsv1.2  https://2wayssl.local:444

    # 3. view certificate and nginx.conf
    $ ls ~/.2wayssl 
    2wayssl.local.server.crt 2wayssl.local.server.key ca.key                   client.crt               client.key
    2wayssl.local.server.csr ca.crt                   ca.srl                   client.csr               nginx.conf
