@@ -1,7 +1,7 @@
#!/bin/sh

if [ "${ENV}" = "development" ]; then
    cd cmd && exec air
    exec air
else
    exec ../main/main
fi