from subprocess import Popen, PIPE
import os
import time
import asyncio
import locale
import sys
from asyncio.subprocess import PIPE
from contextlib import closing

serverArgs = ["./oauth2local",
              "--verbose",
              "--config",
              "test-oauth2-config.yml",
              "serve"]
print(serverArgs)
callbackArgs = ["./oauth2local",
                "--verbose",
                "--config",
                "test-oauth2-config.yml",
                "callback",
                "loc-auth://callback?code=cryptocode&state=none"]
tokenArgs = ["./oauth2local",
             "--config",
             "test-oauth2-config.yml",
             "token"]


async def readline_and_kill(args, sf, cf, tf, ef):
    # start child process
    server = await asyncio.create_subprocess_exec(
        args[0], *args[1:], stdout=sf, stderr=ef)
    time.sleep(3)
    callbackCmd = Popen(callbackArgs,
                        stdout=cf,
                        stderr=ef)
    callbackCmd.wait()
    cf.close()
    time.sleep(3)
    tokenCmd = Popen(callbackArgs,
                     stdout=tf,
                     stderr=ef)
    tf.close()
    ef.close()
    tokenCmd.wait()
    exitcode = 1
    if tokenCmd.returncode == 0:
        print("Success")
        exitCode= 0
    else:
        with open("token.log", "r") as logFile:
            log = logFile.read()
        print("Error", log)
    server.kill()
    await server.wait()  # wait for the child process to exit
    return exitCode


if sys.platform == "win32":
    loop = asyncio.ProactorEventLoop()
    asyncio.set_event_loop(loop)
else:
    loop = asyncio.get_event_loop()

with closing(loop):
    sys.exit(loop.run_until_complete(
        readline_and_kill(
            serverArgs,
            open("server.log", "w"),
            open("callback.log", "w"),
            open("token.log", "w"),
            open("error.log", "w"),
        )))
