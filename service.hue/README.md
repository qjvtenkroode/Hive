# Philips Hue

Service for all Philips Hue assets and configurations within the home automation system. Responsible for controlling and receiving state on all Philips Hue assets.
Currently only works for regular light assets

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
        "identifier": "lamp 1",
        "state": "on",
        "type": "hue light",
        "last_update": "Oct 14 09:29:17"
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


