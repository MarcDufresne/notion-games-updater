import time
from typing import Any

import httpx


class IGDBClientError(Exception):
    def __init__(self, message: str) -> None:
        self.message = message


class IGDB:
    _AUTH_URL = "https://id.twitch.tv/oauth2/token"
    _API_URL = "https://api.igdb.com/v4/"
    _RETRIES = 3

    def __init__(self, client_id: str, client_secret: str) -> None:
        self.client_id = client_id
        self.__client_secret = client_secret
        self.__access_token = None
        self.__access_token_expires_at = 0

    def _refresh_token(self) -> None:
        with httpx.Client() as client:
            try:
                response = client.post(
                    self._AUTH_URL,
                    params={
                        "client_id": self.client_id,
                        "client_secret": self.__client_secret,
                        "grant_type": "client_credentials",
                    },
                )
                response.raise_for_status()
            except httpx.RequestError as re:
                raise IGDBClientError(f"Failed to refresh token: {re}")
            except httpx.HTTPStatusError as he:
                raise IGDBClientError(
                    f"Failed to refresh token [{he.response.status_code}]: {he.response.text}"
                )
            self.__access_token = response.json()["access_token"]

    def _request(self, endpoint: str, query: str) -> httpx.Response:
        if (
            self.__access_token is None
            or self.__access_token_expires_at - 3600 < time.time()
        ):
            self._refresh_token()

        headers = {
            "Client-ID": self.client_id,
            "Authorization": f"Bearer {self.__access_token}",
        }

        with httpx.Client(base_url=self._API_URL, headers=headers) as client:
            client: httpx.Client
            try:
                response = client.post(endpoint, content=query)
            except httpx.RequestError as re:
                raise IGDBClientError(f"Failed to make request: {re}")

            try_count = 0
            while True:
                try:
                    response.raise_for_status()
                    break
                except httpx.HTTPStatusError as e:
                    if e.response.status_code == 429:
                        try_count = 0
                        time.sleep(0.25)
                        continue
                    try_count += 1
                    if try_count > self._RETRIES:
                        raise IGDBClientError(
                            f"Failed to make request [{e.response.status_code}]: {e.response.text}"
                        )

        return response

    def request(
        self, endpoint: str, query: str
    ) -> dict[str, Any] | list[dict[str, Any]]:
        return self._request(endpoint, query).json()

    def list_all(self, endpoint: str, query: str) -> list[dict[str, Any]]:
        response = self._request(endpoint, query)
        data = response.json()

        if not isinstance(data, list):
            raise IGDBClientError(f"Expected a list, got {type(data)}")

        total = int(response.headers["x-count"])
        offset = len(data)

        while offset < total:
            new_query = f"{query} offset {offset};"
            print(new_query)
            response = self._request(endpoint, new_query)
            data.extend(response.json())
            offset += len(response.json())

        return data
