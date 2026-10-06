from dp_population_types_api_sdk_python.models import HTTPHeaders


def test_to_http_headers_converts_all_values() -> None:
    headers = HTTPHeaders(
        access_token="access-token",
        florence_token="florence-token",
    )

    assert headers.to_http_headers() == {
        "Authorization": "Bearer access-token",
        "X-Florence-Token": "florence-token",
    }


def test_to_http_headers_omits_unset_values() -> None:
    headers = HTTPHeaders(access_token="access-token")

    assert headers.to_http_headers() == {
        "Authorization": "Bearer access-token",
    }
