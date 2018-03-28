title: Mistune Syntax Highlighter, MathJax support and centered images
description: Server rendering syntax highlighting with Mistune and other stuff
slug: mistune-parser-syntax-mathjax-centered-images
date: 2015-09-29 23:51:00
tags:
    - dev
    - markdown
    - python

# Introduction

*This article was written for [markdownblog.com](http://markdownblog.com/) and denotes some changes to the platform. These changes don't apply on smallblog.*

As you may (or may not) know, I recently switched from [Misaka](http://misaka.61924.nl/) ([GitHub](https://github.com/FSX/misaka)) to [Mistune](http://mistune.readthedocs.org/en/latest/) ([GitHub](https://github.com/lepture/mistune)) mainly because Mistune is a pure python markdown parser. It means that it is easier to declare new grammars (and new behaviors) than modifying Misaka itself (which, just as a reminder, is a binding for Sundown, markdown engine written in C). I needed to modify the parser's behavior because [carado](http://carado.markdownblog.com/) asked me for a MathJax support (see [my previous update](http://depado.markdownblog.com/2015-09-19-centered-images-and-mathjax-support) for mor information). It took me a really long time to figure out how I could do that using Mistune, because the documentation about adding additionnal behaviors isn't as clear as I expected.

# MathJax Support

## The Javascript library

First of all, I started by including the library in the base template for every blog. For now, the lib is loaded whenever you're reading a blog even if it doesn't use MathJax. Although that behavior may change as I intend to add article-specific settings. For example if you want to enable or disable some parsing options or modify the way some items are rendered. But that will come later.

```html
<script type="text/x-mathjax-config">
    MathJax.Hub.Config({
      tex2jax: {inlineMath: [['$','$']]}
    });
</script>
<script type="text/javascript" src="https://cdn.mathjax.org/mathjax/latest/MathJax.js?config=TeX-AMS-MML_HTMLorMML"></script>
```

As you can see the delimiters for inline math is `$`. It's obviously easier to parse if the inline blocks are defined using `$...$` and blocks are defined using `$$...$$`.

## The mistune renderer

```python
class MathBlockGrammar(mistune.BlockGrammar):
    block_math = re.compile(r"^\$\$(.*?)\$\$", re.DOTALL)
    latex_environment = re.compile(r"^\\begin\{([a-z]*\*?)\}(.*?)\\end\{\1\}", re.DOTALL)


class MathBlockLexer(mistune.BlockLexer):
    default_rules = ['block_math', 'latex_environment'] + mistune.BlockLexer.default_rules

    def __init__(self, rules=None, **kwargs):
        if rules is None:
            rules = MathBlockGrammar()
        super(MathBlockLexer, self).__init__(rules, **kwargs)

    def parse_block_math(self, m):
        """Parse a $$math$$ block"""
        self.tokens.append({
            'type': 'block_math',
            'text': m.group(1)
        })

    def parse_latex_environment(self, m):
        self.tokens.append({
            'type': 'latex_environment',
            'name': m.group(1),
            'text': m.group(2)
        })


class MathInlineGrammar(mistune.InlineGrammar):
    math = re.compile(r"^\$(.+?)\$", re.DOTALL)
    block_math = re.compile(r"^\$\$(.+?)\$\$", re.DOTALL)
    text = re.compile(r'^[\s\S]+?(?=[\\<!\[_*`~$]|https?://| {2,}\n|$)')


class MathInlineLexer(mistune.InlineLexer):
    default_rules = ['block_math', 'math'] + mistune.InlineLexer.default_rules

    def __init__(self, renderer, rules=None, **kwargs):
        if rules is None:
            rules = MathInlineGrammar()
        super(MathInlineLexer, self).__init__(renderer, rules, **kwargs)

    def output_math(self, m):
        return self.renderer.inline_math(m.group(1))

    def output_block_math(self, m):
        return self.renderer.block_math(m.group(1))


class MarkdownWithMath(mistune.Markdown):
    def __init__(self, renderer, **kwargs):
        if 'inline' not in kwargs:
            kwargs['inline'] = MathInlineLexer
        if 'block' not in kwargs:
            kwargs['block'] = MathBlockLexer
        super(MarkdownWithMath, self).__init__(renderer, **kwargs)

    def output_block_math(self):
        return self.renderer.block_math(self.token['text'])

    def output_latex_environment(self):
        return self.renderer.latex_environment(self.token['name'], self.token['text'])
```

Let's start by wondering what we're trying to achieve here. MathJax is a javascript library that will read the content of the page it is included in and replace everything that is between `$` or `$$` delimiters with custom html/css (or even svg) to display beautiful math. This means that the only thing we're trying to achieve there is just to **not parse** what's between those delimiters.   

Now this is getting a bit complicated here. Mainly due to how Mistune works, you first need to define a grammar class (`MathInlineGrammar` and `MathBlockGrammar`) which mainly consist of regular expressions for the things to match in the markdown text. We then declare the lexers that will use the grammar defined earlier (`MathInlineLexer` and `MathBlockLexer`)

We'll then define a new markdown parser that will be able to understand and use the two lexers defined above (`MarkdownWithMath`).

# Syntax Highlighter

This part is quite similar to the article I wrote about how to do a [markdown syntax highlighter with Misaka](http://depado.markdownblog.com/2015-02-05-markdown-syntax-highlighter). The code doesn't change much.

```python
class HighlighterRenderer(mistune.Renderer):

    def block_code(self, code, lang=None):
        if not lang:
            lang = 'text'
        try:
            lexer = get_lexer_by_name(lang, stripall=True)
        except:
            lexer = get_lexer_by_name('text', stripall=True)
        formatter = HtmlFormatter()
        return "{open_block}{formatted}{close_block}".format(
            open_block="<div class='code-highlight'>" if lang != 'text' else '',
            formatted=highlight(code, lexer, formatter),
            close_block="</div>" if lang != 'text' else ''
        )
```

Now let's add the support for MathJax and centered images in our renderer :

```python
class HighlighterRenderer(mistune.Renderer):

    def block_code(self, code, lang=None):
        if not lang:
            lang = 'text'
        try:
            lexer = get_lexer_by_name(lang, stripall=True)
        except:
            lexer = get_lexer_by_name('text', stripall=True)
        formatter = HtmlFormatter()
        return "{open_block}{formatted}{close_block}".format(
            open_block="<div class='code-highlight'>" if lang != 'text' else '',
            formatted=highlight(code, lexer, formatter),
            close_block="</div>" if lang != 'text' else ''
        )

    def table(self, header, body):
        return "<table class='table table-bordered table-hover'>" + header + body + "</table>"

    def image(self, src, title, text):
        if src.startswith('javascript:'):
            src = ''
        text = mistune.escape(text, quote=True)
        if title:
            title = mistune.escape(title, quote=True)
            html = '<img class="img-responsive center-block" src="%s" alt="%s" title="%s"' % (src, text, title)
        else:
            html = '<img class="img-responsive center-block" src="%s" alt="%s"' % (src, text)
        if self.options.get('use_xhtml'):
            return '%s />' % html
        return '%s>' % html

    # Pass math through unaltered - mathjax does the rendering in the browser
    def block_math(self, text):
        return '$$%s$$' % text

    def latex_environment(self, name, text):
        return r'\begin{%s}%s\end{%s}' % (name, text, name)

    def inline_math(self, text):
        return '$%s$' % text
```
Now that we are all set, we can initialize the renderer like this :

```python
markdown_renderer = MarkdownWithMath(renderer=HighlighterRenderer(escape=False))
```
