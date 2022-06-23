# vacadm

This repository was created as part of my final apprenticeship project.
It is a REST interface for vacation management. Details on how to use it can be accessed via a Swagger UI. A library overview is located [here (pkg.go.dev)](https://pkg.go.dev/github.com/miyunari/vacadm).
The final version is published as [release 1.0.0](https://github.com/miyunari/vacadm/releases/tag/v1.0.0). Other assets such as the project application, project documentation and project presentation are also located there. 

## Usage

```bash
Usage of ./vacadm:
  -address string
    	ip:port (default "localhost:8080")
  -init.root
    	create root user on startup
  -secret string
    	secret for jwt token
  -smtp.host string
    	address of smtp server
  -smtp.password string
    	smtp user password
  -smtp.port string
    	port of smtp server
  -smtp.user string
    	smtp user mail address
  -sql.conn string
    	sql connection str. user:password@/dbname
    			example: root:my-secret-pw@(127.0.0.1:3306)/test?parseTime=true
  -swagger.enable
    	enables /swagger endpoint
  -timeout duration
    	server timeout (default 1m0s)
```