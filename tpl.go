package main

const NGX_TEMPLATE = `
server {
	  listen 444; # change to 443 in production
	  server_name {DOMAIN};
	  ssl on;
	  #server auth
	  ssl_certificate {SERVER_CRT_PATH};
	  ssl_certificate_key {SERVER_KEY_PATH};

	  #client auth
	  ssl_verify_client on; 
	  ssl_client_certificate {CA_CRT_PATH};

	  #ssl options
	  ssl_protocols SSLv2 SSLv3 TLSv1 TLSv1.1 TLSv1.2; 
	  ssl_ciphers 'ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-AES256-GCM-SHA384:kEDH+AESGCM:ECDHE-RSA-AES128-SHA256:ECDHE-ECDSA-AES128-SHA256:ECDHE-RSA-AES128-SHA:ECDHE-ECDSA-AES128-SHA:ECDHE-RSA-AES256-SHA384:ECDHE-ECDSA-AES256-SHA384:ECDHE-RSA-AES256-SHA:ECDHE-ECDSA-AES256-SHA:DHE-RSA-AES128-SHA256:DHE-RSA-AES128-SHA:DHE-RSA-AES256-SHA256:DHE-DSS-AES256-SHA:AES128-GCM-SHA256:AES256-GCM-SHA384:ECDHE-RSA-RC4-SHA:ECDHE-ECDSA-RC4-SHA:RC4-SHA:HIGH:!aNULL:!eNULL:!EXPORT:!DES:!3DES:!MD5:!PSK';
	  ssl_prefer_server_ciphers  on;
	  ssl_session_timeout 5m;

	  location / {
			  proxy_pass http://127.0.0.1:4500;
			  proxy_set_header Host $host:$server_port;
			  proxy_set_header X-Real-IP $remote_addr;
			  proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
			  proxy_next_upstream http_502 http_504 error timeout invalid_header; 
	  }
}
`
