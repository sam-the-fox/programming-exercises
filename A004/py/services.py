import math
import threading
import time
import tomllib

import requests


class Service:
    name: str
    url: str
    up: bool
    latency_ms: int

    def __init__(self, name: str, url: str):
        self.name = name
        self.url = url
        self.up = False
        self.latency_ms = 0


class Services(threading.Thread):
    all: list[Service] = []  # noqa: RUF012

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)

        with open("services.toml", "r") as f:
            config = tomllib.loads(f.read())
            for service in config["services"]:
                Services.all.append(Service(service["name"], service["url"]))

    def run(self, *args, **kwargs):
        while True:
            print(f"meow. querying {len(Services.all)} services.")
            for service in Services.all:
                print(f"Querying {service.url}")

                start = time.time()
                try:
                    resp = requests.get(service.url)
                    service.up = resp.status_code == 200
                    service.latency_ms = int((time.time() - start) * 1000)
                except:  # noqa: E722
                    service.up = False
                    service.latency_ms = int(math.pow(2, 31)) - 1

                print(f"Queried {service.url} in {service.latency_ms}ms")

            time.sleep(60)


s = Services()
s.start()
