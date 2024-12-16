# Design

1. Ask user for first article
2. Ask user for second article
3. Make sure articles exist
   1. Call API (It does some nice filtering)
4. Check Redis Cache
   1. If in cache then we are done
5. If articles don't exits then reprompt user for article names
6. If articles do exist then start breadth first search until a path is found.

## Database Details

Titles should be the last part of the href path. For example Adolf Hitlers database title is Adolf_Hitler and can be found at `/html/head/link[2]`

## Search Algorithm Details

struct node {article title, child_nodes}

func breadthFirstSearch() -> path

1. Make request channel (setup function)
2. Make parent node (include links)
3. Call search (node, channel, depth, maxdepth) --Keep increasing maxdepth until a solution is found, or max depth is reached

fun search(node, channel, depth, maxdepth) -> path --Behavior undefined for depth limit that is too large

1. Recurse until depth == maxDepth
2. For each link (Concurrent)
   1. Query database for it, else send it to the channel
      1. If not in database add it
   2. Check if any of the children are the target
      1. if so return
      2. else return empty string

### Rate limited Function details

```
struct request (title, priority, return channel) // priority should be equal to depth, lower means it needs to be done first. Priority -1 means its time to stop.

func apiRequestProcessor(bufferedChannel(request))
    queue of requests
    Every second:
        pop 190 requests off queue and make go routines to call them
        empty buffered channel into queue
```

## HTML parsing notes

* Look at 'href="/wiki/'
* Do not look at href="/wiki/File
* Be sure to filter out duplicates
* The current article title can be found in the content-location header
