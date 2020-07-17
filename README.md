# srcf.redact

This Mattermost plugin lets team administrators delete old messages from their
channels based on their age. It provides two commands:

 * `/channel_id` returns the channel id
 * `/delete_posts` this is the command that actually deletes the messages.

The delete command is currently only able to delete messages older than a
specific number of days. Pull requests to enhance its functionality (e.g. to
allow messages sent between two dates, or other duration formats) are welcome.
