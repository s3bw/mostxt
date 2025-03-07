<div align="center">
    <h1>Mostxt</h1>
</div>

<div align="center">
    From `template` -> `cli`
</div>


## Templates

```
Title: {{ title example('My Blog Title') describe('Make it catchy') }}
Summary: {{ summary describe('A sort summary of the text') }}
Date: {{ x:datetime 'YYYY' }}
Tags: {{ tags:list }}
```

## CLI

```
$ mostxt template.md output.md

Enter title (e.g: My Blog Title)
Make it catchy
$ How to lose a guy in 10 days

Enter summary
A sort summary of the text
$ A tale about falling inlove

Enter tags
$ interesting
$ catchy
$ truestory
$

Output written to output.md
```
