import subprocess


def getAccessToken() -> str:
    res = subprocess.check_output(["oauth2local", "token"])
    return res.decode("utf-8").strip()


if __name__ == '__main__':
    print(getAccessToken())
