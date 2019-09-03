# Registry

Registery service for all assets and configurations within the home automation system. Responsible for keeping a record of all active assets and distributing asset/system configuration.

## Endpoints

**Definition**

`GET /assets`

**Description**

List all assets known in the registry.

**Arguments**

None

**Response**

- `200 OK` on success

    ```json
    [
        {
            "identifier": "asdagsdfj",
            "name": "Switch bedroom",
            "type": "Shelly 1"
        },
        {
            "identifier": "aogiuhggi",
            "name": "Switch living room",
            "type": "Shelly 1"
        }
    ]
    ```

---

**Definition**
 
`GET /assets/<identifier>`

**Description**

Get all stored information about a specific asset

**Arguments**

None

**Response**

- `200 OK` on success

    ```json
    {
        "identifier": "asdagsdfj",
        "name": "Switch bedroom",
        "type": "Shelly 1"
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

`POST /assets`

**Description**

Register a new asset

**Arguments**

- `"identifier":string` a globally unique identifier for this asset
- `"name":string` a friendly name for this asset
- `"type":string` the type of asset

**Response**

- `201 Created` on success

    ```json
    {
        "status": "ok",
        "message": "asset created"
    }
    ```

---

**Definition**

`DELETE /assets/<identifier>`

**Description**

Delete an asset

**Arguments**

None

**Response**

- `204 No Content` on success

    ```json
    {
        "status": "ok",
        "message": "asset deleted succesfully"
    }
    ```

- `404 Not Found` if asset is not found

    ```json
    {
        "status": "error",
        "message": "asset not found"
    }
    ```

