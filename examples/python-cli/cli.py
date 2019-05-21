import subprocess


def getAccessToken() -> str:
    res = subprocess.run(["oauth2local", "token"], stdout=subprocess.PIPE)
    return res.stdout.decode("utf-8").strip()


if __name__ == '__main__':
    print(getAccessToken())
