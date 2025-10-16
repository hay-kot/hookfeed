-- Example transform function that extracts and reformats webhook data
function transform(input)
    local output = {
        transformed = true,
        original_event = input.body.event,
        processed_at = os.date("!%Y-%m-%dT%H:%M:%SZ"),
        payload = {
            id = input.body.data.id,
            status = input.body.data.status,
            received_timestamp = input.body.timestamp
        }
    }

    return output
end
