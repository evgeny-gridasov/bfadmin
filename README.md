### BFADMIN

BFADMIN is a Web-based administration tool to control a game server. It may be used to
change configuration settings, select maps for rotation, startup and shutdown server processes.

Games supported by BFADMIN:

- Battlefield Vietnam
- Battlefield 1942
- Unreal Tournament 2004

It is recommended to put BFADMIN behind NGINX and password protect it.

I wrote this tool to help a group of friends start/stop/reconfigure BFV and BF1942 dedicated servers without the need to ssh onto the server every time.
As an exercise, I used Golang. It is far from being perfect, has a few TODOs and may be bugs, but does the job! 
