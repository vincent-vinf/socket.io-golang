echo off
set func=%1
set version=%2

IF %func%==public (
    echo Start
    git tag %version%

    git push origin %version%

    SET GOPROXY=proxy.golang.org

    go list -m github.com/doquangtan/gofiber-socket.io@%version%
    
    echo Done
)

@REM make public v0.1.1