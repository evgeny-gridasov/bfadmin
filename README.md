## BFADMIN - Game Admin Console

### Overview

BFADMIN is a Web-based administration console for managing game servers. It may be used to
change configuration settings, select maps for rotation, startup and shutdown server processes.

Games supported by BFADMIN:

- Battlefield Vietnam
- Battlefield 1942
- Battlefield 2
- Unreal Tournament 2004

It is recommended to put BFADMIN behind NGINX and password protect it.

### Building

There are no external dependencies, and BFADMIN should build without any issues. Just run:
```
go build
```

Copy the supplied `bfadmin.conf.default` configuration file to `bfadmin.conf` and change configuration parameters to your requirements. The paramter names should be self-explanatory.

### About

I wrote BFADMIN to help a group of friends start/stop/reconfigure BFV and BF1942 dedicated servers without the need to ssh onto the server every time.
As an exercise, I used Golang. It is far from being perfect, has a few TODOs and may be bugs, but does the job! 

Written by Evgeny Gridasov ([@evgeny-gridasov](http://github.com/evgeny-gridasov))
