echo off
set func=%1
set version=%2

IF %func%==public (
    echo Start %func%:
    git checkout main

    git pull

    git tag %version%

    git push origin %version%

    SET GOPROXY=proxy.golang.org

    go list -m github.com/doquangtan/gofiber-socket.io@%version%
    
    echo Done %func%
)

@REM make public v0.1.8