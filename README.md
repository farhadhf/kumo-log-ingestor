A simple golang program for ingesting [KumoMTA](https://kumomta.com) log webhook and inserting them into a PostgreSQL database.


### Run the program
```bash
$ export DATABASE_URL='postgresql://username:password@host:port/db_name?sslmode=disable'
$ export LISTEN_ADDR='localhost:3000'
$ go build main.go
$ ./main
```

### Set up KumoMTA log webhook
Add the following code in `init.lua`, before setting up the queue helper. Learn more about [KumoMTA webhooks](https://docs.kumomta.com/userguide/operation/webhooks/).
```lua
-- Replace with the actual address of the postgres ingestion program
local LOG_WEBHOOK_URL = 'http://localhost:3000'

-- IMPORTANT: This needs to be before defining the queue helper.
log_hooks:new_json {
  name = "postgres-webhook",
  -- log_parameters are combined with the name and
  -- passed through to kumo.configure_log_hook
  log_parameters = {
    -- Replace with the headers you want to have in postgres
    headers = { 'Subject', 'From', 'Message-Id', 'References', 'In-Reply-To' },
    -- Add meta that you want to have in postgres
    meta = { 'tenant', 'domain_id', 'customer_id', 'direction', 'metadata_retention', 'data_retention', 'tags' }
  },
  -- queue config are passed to kumo.make_queue_config.
  -- You can use these to override the retry parameters
  -- if you wish.
  -- The defaults are shown below.
  queue_config = {
    retry_interval = "1m",
    max_retry_interval = "20m",
  },
  -- The URL to POST the JSON to
  url = LOG_WEBHOOK_URL
}
```
