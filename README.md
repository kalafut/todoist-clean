# Todoist De-Thingser

If you're migrating to Todoist from Things 3 or use both apps, copy/paste from Things
to Todoist works pretty well. You can easily select tasks in bulk in Things, and paste
into Todoist allows you the choice to make multiple tasks. It works fine except the new
tasks in Todoist look like:

```
[ ] Pick up milk
[ ] Take out trash
[ ] Make doctor appt
```

This small Go app will update all active tasks to remove that leading `[ ]` visual clutter.

## Requirements

I've not packaged this up as a binary, so you'll need to have Go installed to build or run it.
You can get Go here: https://go.dev/doc/install

## Usage

1. Install the app: `go install github.com/kalafut/todoist-clean@latest`
1. Get an API token from Todoist: https://todoist.com/app/settings/integrations/developer
1. Set the `TODOIST_API_TOKEN` environment variable to the token you got above: `export TODOIST_API_TOKEN=<your token>`

1. Run `todoist-clean` (it will show you what it wants to change before making the changes)

Note that the Todoist API has a rate limit of 450 requests per 15 minutes. If you hit a rate limit (unlikely
unless you have a huge number of tasks to fix), just rerun in a few minutes and the updates will resume.
