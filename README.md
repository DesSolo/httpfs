# httpfs
Простое файловое хранилище поверх http

## Usage
``` bash
# upload
$ curl -T file.txt http://127.0.0.1:8080
{"hash":"47a013e660d408619d894b20806b1d5086aab03b"}

# download
$ curl http://127.0.0.1:8080/47a013e660d408619d894b20806b1d5086aab03b
Hello world!

# actual fs file
$ cat /tmp/httpfs/47/47a013e660d408619d894b20806b1d5086aab03b 
Hello world!

# sha1 sum
$ sha1sum /tmp/httpfs/47/47a013e660d408619d894b20806b1d5086aab03b
47a013e660d408619d894b20806b1d5086aab03b  /tmp/httpfs/47/47a013e660d408619d894b20806b1d5086aab03b
```