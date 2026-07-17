# 1) Initial request advertising a client-side tool
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "stream": true,
    "messages": [{"role": "user", "content": "What time is it?"}],
    "tools": [{
      "type": "function",
      "function": {
        "name": "get_current_time",
        "description": "Get the current time",
        "parameters": {"type": "object", "properties": {}}
      }
    }]
  }'

# 2) Stream ends with finish_reason: "tool_calls". The client executes the tool
#    locally and submits the result back in a follow-up request:
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "stream": true,
    "tools": [ ... same tools array ... ],
    "messages": [
      {"role": "user", "content": "What time is it?"},
      {"role": "assistant", "content": null, "tool_calls": [
        {"id": "call_abc", "type": "function",
         "function": {"name": "get_current_time", "arguments": "{}"}}
      ]},
      {"role": "tool", "tool_call_id": "call_abc", "content": "2026-04-30T14:30:00Z"}
    ]
  }'
