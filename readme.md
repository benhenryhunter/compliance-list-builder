# Compliance List Builder

This a pretty basic compliance list builder that subscribes to a websocket eth client and tracks any Tornado Cash Deposit events. When an event is found this tool will track the from address for the event. 

The every 5 minutes the latest addresses are written to the addresses.json file, or another file if specified.
