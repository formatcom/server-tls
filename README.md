REF: https://docs.docker.com/engine/security/https/

 * OpenSSL
 * x509
 * TLS

1.- Autoridad Certificadora

Generate a certificate authority private key (CA)
~~~
$ openssl genrsa -aes256 -out ca.key.pem 4096
~~~

Generate the root CA certificate (CSR)
~~~
$ openssl req -new -key ca.key.pem -x509 -days 365 -sha256 -out ca.pem
~~~

Ver resultado
~~~
$ openssl x509 -in ca.pem -text -noout | less
~~~

2.- Generar la clave del servidor
~~~
$ openssl genrsa -out server.key.pem 4096
~~~

Solicitud de firma de certificado (CSR) Asegúrese de que el "Nombre Común"
coincida con el nombre del host que utiliza para conectarse al demonio. 
~~~
$ openssl req -subj "/CN=$HOST" -sha256 -new -key server.key.pem -out server.csr
~~~

Dado que las conexiones TLS pueden realizarse a través de la dirección IP y el nombre DNS,
las direcciones IP deben especificarse al crear el certificado.
Por ejemplo, para permitir conexiones usando 192.168.0.10 y 127.0.0.1

~~~
$ echo subjectAltName = DNS:$HOST,IP:192.168.0.10,IP:127.0.0.1 >> extfile.cnf
~~~

Configure los atributos de uso extendido de la clave del demonio Docker para que solo
se usen para la autenticación del servidor

~~~
$ echo extendedKeyUsage = serverAuth >> extfile.cnf
~~~

Validar que la licencia del cliente no se ha modificado desde que se genero

~~~
$ echo keyUsage = digitalSignature >> extfile.cnf
~~~

Firmar la clave publica.
~~~
$ openssl x509 -req -days 365 -sha256 -in server.csr -CA ca.pem -CAkey ca.key.pem \
  -CAcreateserial -out server.pem -extfile extfile.cnf
~~~

Ver resultado
~~~
$ openssl x509 -in server.pem -text -noout | less
~~~

3.- Generar Cliente

Certificate authority (CA)
~~~
$ openssl genrsa -out client.key.pem 4096
~~~

Certificate signing request (CSR)
~~~
$ openssl req -subj "/CN=Vinicio Valbuena" -new -key client.key.pem -out client.csr
~~~

Para que la clave sea adecuada para la autenticación del cliente, cree un nuevo archivo
de configuración de extensiones

~~~
$ echo extendedKeyUsage = clientAuth > extfile-client.cnf
~~~

Generar el certificado firmado
~~~
$ openssl x509 -req -days 365 -sha256 -in client.csr -CA ca.pem -CAkey ca.key.pem \
  -CAcreateserial -out client.pem -extfile extfile-client.cnf
~~~

Ver resultado
~~~
$ openssl x509 -in client.pem -text -noout | less
~~~

Limpiar lo que no sea necesario
~~~
$ rm -v client.csr server.csr extfile.cnf extfile-client.cnf
~~~

Agragar seguridad a las claves
~~~
$ chmod -v 0400 ca.key.pem client.key.pem server.key.pem
$ chmod -v 0444 ca.pem server.pem client.pem
~~~

curl -vvv --cert client.pem --key client.key.pem https://0.0.0.0:1200 -k -d 'hola mundo'

