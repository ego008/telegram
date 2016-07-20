package main

const startMsgTpl = `Hi %s!

This is the official @HentaiDB bot. You can browse the Danbooru pics, GIF's and videos here.

You can also use it to search and share content with your friends.
Just type "@HentaiDBot hatsune_miku" in any chat and select the result you want to send.`

const helpMsg = `<b>tag1 tag2</b>
<code>Search for posts that have tag1 and tag2.</code>

<b>~tag1 ~tag2</b>
<code>Search for posts that have tag1 or tag2. (Currently does not work)</code>

<b>night~</b>
<code>Fuzzy search for the tag night. This will return results such as night fight bright and so on according to the </code><a href=\"https://en.wikipedia.org/wiki/Levenshtein_distance\">Levenshtein distance</a><code>.</code>

<b>-tag1</b>
<code>Search for posts that don't have tag1.</code>

<b>ta*1</b>
<code>Search for posts with tags that starts with ta and ends with 1.</code>

<b>user:bob</b>
<code>Search for posts uploaded by the user Bob.</code>

<b>md5:foo</b>
<code>Search for posts with the MD5 hash foo.</code>

<b>md5:foo*</b>
<code>Search for posts whose MD5 starts with the MD5 hash foo.</code>

<b>rating:questionable</b>
<code>Search for posts that are rated questionable.</code>

<b>-rating:questionable</b>
<code>Search for posts that are not rated questionable.</code>

<b>parent:1234</b>
<code>Search for posts that have 1234 as a parent (and include post 1234).</code>

<b>rating:questionable rating:safe</b>
<code>In general, combining the same metatags (the ones that have colons in them) will not work.</code>

<b>rating:questionable parent:100</b>
<code>You can combine different metatags, however.</code>

<b>width:&gt;=1000 height:&gt;1000</b>
<code>Find images with a width greater than or equal to 1000 and a height greater than 1000.</code>

<b>score:&gt;=10</b>
<code>Find images with a score greater than or equal to 10. This value is updated once daily at 12AM CST.</code>

<b>sort:updated:desc</b>
<code>Sort posts by their most recently updated order.</code>

<b>Other sortable types:</b>
<code>- id
- score
- rating
- user
- height
- width
- parent
- source
- updated
Can be sorted by both asc or desc.</code>`
