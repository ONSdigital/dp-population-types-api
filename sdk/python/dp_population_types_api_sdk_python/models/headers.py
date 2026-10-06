from pydantic import BaseModel


class HTTPHeaders(BaseModel):
    access_token: str | None = None
    florence_token: str | None = None

    def to_http_headers(self) -> dict[str, str]:
        headers: dict[str, str] = {}

        if self.access_token is not None:
            headers["Authorization"] = f"Bearer {self.access_token}"
        if self.florence_token is not None:
            headers["X-Florence-Token"] = self.florence_token

        return headers
