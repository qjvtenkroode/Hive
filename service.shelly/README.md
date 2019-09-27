# Shelly

Service for all Shelly assets and configurations within the home automation system. Responsible for controlling and receiving state on all Shelly assets.
Currently only works for shelly1 assets

## Endpoints

**Definition**

`GET /state/<identifier>`

**Description**

List the state of all relays for the assets known by <identifier>.

**Arguments**

None

**Response**

- `200 OK` on success

    ```json
    {
        "identifier": "asdagsdfj",
        "state": "on"
    }
    ```

- `404 Not Found` if asset is not found

    ```json
    {
        "status": "error",
        "message": "asset not found"
    }
    ```

---

**Definition**

`POST /state/<identifier>`

**Description**

Toggle the state

**Arguments**

```json
{
    "identifier": "asdagsdfj",
    "state": "off"
}
```


**Response**

- `200 OK` on success

    ```json
    {
        "identifier": "asdagsdfj",
        "state": "off"
    }
    ```

- `404 Not Found` if asset is not found

    ```json
    {
        "status": "error",
        "message": "asset not found"
    }
    ```

---
