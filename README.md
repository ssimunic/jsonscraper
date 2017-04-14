# jsonscraper

JSON configurable concurrent scraper. Written in Go.

For given JSON config file(s), produces JSON file(s) with results.

## Instructions

`$ jsonscraper configPath`

Output will be saved in file path that is provided in configuration.

You can also run it with multiple configs at once: 
`$ jsonscraper configPath1 configPath2 configPath3 ...`

## Config documentation
Configuration is an object which consists of `urls`, `targets` and `output`.

### `urls` (Array)
Array of URLs which will be scraped.

### `targets` (Array)
Each object in this array should have atleast:
* `selector` is selection query to be perfomed
* `type` can be inner HTML (`html`), plain text (`text`) or attribute (for example `attr:href`) 
* `tag` is field key which appears in output file

Optional is `submatch` which can contain regular expression. `selector` can be omitted if `submatch` is present, then whole document will be used for lookups.

### `output` (Object)
* `path` is an output path where data will be saved. `$FILENAME` will be replaced with input file name.


## Example

#### Input file: `input/example.json`

```json
{
    "urls": [
        "https://news.ycombinator.com/"
    ],
    "targets": [
        {
            "selector": ".storylink",
            "type": "text",
            "tag": "storyTitleText"
        },
        {
            "selector": ".title",
            "type": "html",
            "tag": "storyTitleHtml"
        },
        {
            "selector": ".storylink",
            "type": "text",
            "tag": "storyTitleWords",
            "submatch": "([a-zA-Z]+)+"
        },
                {
            "selector": ".storylink",
            "type": "attr:href",
            "tag": "storyTitleLinks"
        }
    ],
    "output": {
        "path": "output/$FILENAME"
    }
}
```

Above configuration will produce following data:

#### Output file: `output/example.json`

```json
{
   "storyTitleHtml":[
      "<span class=\"rank\">1.</span>",
      "<a href=\"https://redditblog.com/2017/04/13/how-we-built-rplace/\" class=\"storylink\">How We Built r/Place</a><span class=\"sitebit comhead\"> (<a href=\"from?site=redditblog.com\"><span class=\"sitestr\">redditblog.com</span></a>)</span>",
      "<span class=\"rank\">2.</span>",
      ...
   ],
   "storyTitleLinks":[
      "https://redditblog.com/2017/04/13/how-we-built-rplace/",
      "http://www.bbc.com/news/science-environment-39592059",
      "https://stripe.com/blog/increment",
      ...
   ],
   "storyTitleText":[
      "How We Built r/Place",
      "Saturn moon 'able to support life'",
      "Introducing Increment",
      ...
   ],
   "storyTitleWords":[
      "How",
      "We",
      "Built",
      ...
   ]
}
```